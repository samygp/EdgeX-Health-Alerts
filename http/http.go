package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/labstack/echo"
	"github.com/samygp/edgex-health-alerts/fault"
	"github.com/samygp/edgex-health-alerts/log"
	"github.com/samygp/edgex-health-alerts/version"
)

const (
	// StatusOK represents a http OK status.
	StatusOK = http.StatusOK
	// StatusCreated represents a http Created status.
	StatusCreated = http.StatusCreated
	// StatusNoContent represents a http NoContent status.
	StatusNoContent = http.StatusNoContent
)

const (
	// UnexpectedStatusCode is a fault code for an unexpected status code.
	UnexpectedStatusCode fault.Code = "unexpected_status_code"
)

var statusCodeToStatus = map[int]fault.Status{
	http.StatusNotFound:     fault.NotFound,
	http.StatusForbidden:    fault.PermissionDenied,
	http.StatusBadRequest:   fault.InvalidArgument,
	http.StatusUnauthorized: fault.Unauthenticated,
}

type option struct {
	url        string
	method     string
	body       interface{}
	statusCode int
	response   interface{}
	headers    map[string]string
}

func newOption(url, method string, opts ...Option) *option {
	options := &option{
		url:        url,
		method:     method,
		statusCode: http.StatusOK,
		headers:    make(map[string]string),
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

// Option represents a request option.
type Option func(*option)

// WithStatusCode sets the expected status code of response.
func WithStatusCode(statusCode int) Option {
	return func(opts *option) {
		opts.statusCode = statusCode
	}
}

// WithBody sets the body of the request.
func WithBody(body interface{}) Option {
	return func(opts *option) {
		opts.body = body
	}
}

// WithResponse sets the body of the response.
func WithResponse(out interface{}) Option {
	return func(opts *option) {
		opts.response = out
	}
}

// WithHeader sets a header for the request.
func WithHeader(key, value string) Option {
	return func(opts *option) {
		opts.headers[key] = value
	}
}

// AuthenticateRequest is a function that authenticates a given request.
type AuthenticateRequest func(req *http.Request, fresh bool) error

type errPayload struct {
	Status string `json:"status"`
	Error  struct {
		Detail  string     `json:"detail"`
		Message fault.Code `json:"message"`
	} `json:"error"`
}

// Client is used to make HTTP requests.
type Client interface {
	BaseURL() string

	Get(ctx context.Context, url string, opts ...Option) error
	Post(ctx context.Context, url string, opts ...Option) error
	Patch(ctx context.Context, url string, opts ...Option) error
	Delete(ctx context.Context, url string, opts ...Option) error
}

type client struct {
	baseURL     string
	client      *retryablehttp.Client
	authRequest AuthenticateRequest
}

// New instantiates the Client implementation.
func New(baseURL string, authRequest AuthenticateRequest, debug bool) Client {
	c := &client{
		baseURL: baseURL,
		client:  retryablehttp.NewClient(),
	}

	c.authRequest = authRequest
	c.client.Logger = nil
	c.client.CheckRetry = retryPolicy(c)

	return c
}

func retryPolicy(c *client) retryablehttp.CheckRetry {
	return func(ctx context.Context, res *http.Response, err error) (bool, error) {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		if err != nil {
			return true, err
		}

		if res.StatusCode == 0 || (res.StatusCode == 429) || (res.StatusCode >= 500 && res.StatusCode != 501) {
			return true, nil
		}

		return false, nil
	}
}

func (c *client) BaseURL() string {
	return c.baseURL
}

func (c *client) Get(ctx context.Context, url string, opts ...Option) error {
	return c.do(ctx, url, http.MethodGet, opts...)
}

func (c *client) Post(ctx context.Context, url string, opts ...Option) error {
	return c.do(ctx, url, http.MethodPost, opts...)
}

func (c *client) Patch(ctx context.Context, url string, opts ...Option) error {
	return c.do(ctx, url, http.MethodPatch, opts...)
}

func (c *client) Delete(ctx context.Context, url string, opts ...Option) error {
	return c.do(ctx, url, http.MethodDelete, opts...)
}

func (c *client) do(ctx context.Context, url, method string, opts ...Option) error {
	if ctx == nil {
		return errors.New("The provided ctx must be non-nil")
	}

	opt := newOption(url, method, opts...)
	var reader io.ReadSeeker
	if opt.body != nil {
		buf, err := json.Marshal(opt.body)
		if err != nil {
			return err
		}

		reader = bytes.NewReader(buf)
	}

	req, err := retryablehttp.NewRequest(method, fmt.Sprintf("%s/%s", c.baseURL, url), reader)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	if id := ContextID(ctx); id != "" {
		req.Header.Set(echo.HeaderXRequestID, id)
	}

	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", version.Name, version.Version))

	if opt.body != nil {
		req.Header.Set(echo.HeaderContentType, "application/json")
	}

	for k, v := range opt.headers {
		req.Header.Set(k, v)
	}

	if c.authRequest != nil {
		err := c.authRequest(req.Request, false)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Processing %s request to: %s\n", req.Method, req.URL)
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			logger := ContextLogger(res.Request.Context())
			logger.Errorf("Unable to close response body: %v", err)
		}
	}()

	if opt.response != nil {
		if err := json.NewDecoder(reader).Decode(opt.response); err != nil {
			if responseData, err := ioutil.ReadAll(res.Body); err != nil {
				log.Logger.Errorf("Error parsing response body: %s", err.Error())
			} else {
				opt.response = string(responseData)
				log.Logger.Debugf("Response: %s", opt.response)
			}
		} else {
			log.Logger.Debugf("Response: %#v", opt.response)
		}
	}

	if res.StatusCode != opt.statusCode {
		if s, ok := statusCodeToStatus[res.StatusCode]; ok {
			return fault.New(s, UnexpectedStatusCode, "Unexpected status code %d", res.StatusCode)
		}
		return fault.New(fault.Unknown, UnexpectedStatusCode, "Unexpected status code %d", res.StatusCode)
	}

	return nil
}
