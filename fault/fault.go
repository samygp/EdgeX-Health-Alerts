package fault

import (
	"fmt"
)

// Code represents a fault code.
type Code string

// Status represetns a fault status.
type Status string

const (
	// Canceled represents a RequestTimeout fault status.
	Canceled Status = "canceled"
	// Unknown represents a InternalServerError fault status.
	Unknown Status = "unknown"
	// InvalidArgument represents a BadRequest fault status.
	InvalidArgument Status = "invalid_argument"
	// DeadlineExceeded represents a GatewayTimeout fault status.
	DeadlineExceeded Status = "deadline_exceeded"
	// NotFound represents a NotFound fault status.
	NotFound Status = "not_found"
	// Conflict represents a Conflict fault status.
	Conflict Status = "conflict"
	// PermissionDenied represents a Forbidden fault status.
	PermissionDenied Status = "permission_denied"
	// ResourceExhausted represents a TooManyRequests fault status.
	ResourceExhausted Status = "resource_exhausted"
	// FailedPrecondition represents a PreconditionFailed fault status.
	FailedPrecondition Status = "failed_precondition"
	// Aborted represents a Conflict fault status.
	Aborted Status = "aborted"
	// OutOfRange represents a BadRequest fault status.
	OutOfRange Status = "out_of_range"
	// Unimplemented represents a NotImplemented fault status.
	Unimplemented Status = "unimplemented"
	// Internal represents a InternalServerError fault status.
	Internal Status = "internal"
	// Unavailable represents a InternalServerError fault status.
	Unavailable Status = "unavailable"
	// DataLoss represents a ServiceUnavailable fault status.
	DataLoss Status = "data_loss"
	// Unauthenticated represents a Unauthorized fault status.
	Unauthenticated Status = "unauthenticated"
)

// Fault wraps lower level errors with status, code, message and an original error.
type Fault interface {
	error

	Status() Status
	Code() Code
	Message() string
	OrigErr() error
}

type fault struct {
	status  Status
	code    Code
	message string
	errs    []error
}

// NewFrom returns an Error object described by the status, code, message and origErr.
func NewFrom(origErr error, status Status, code Code, format string, args ...interface{}) Fault {
	var errs []error
	if origErr != nil {
		errs = append(errs, origErr)
	}

	return newFault(status, code, fmt.Sprintf(format, args...), errs)
}

// New returns an Error object described by the status, code and message.
func New(status Status, code Code, format string, args ...interface{}) Fault {
	return newFault(status, code, fmt.Sprintf(format, args...), nil)
}

// Error returns the string representation of the error.
func (f *fault) Error() string {
	msg := fmt.Sprintf("%s:", f.status)

	if len(f.code) > 0 {
		msg = fmt.Sprintf("%s %s -", msg, f.code)
	}

	msg = fmt.Sprintf("%s %s", msg, f.message)

	if len(f.errs) > 0 {
		msg = fmt.Sprintf("%s\ncaused by: %s", msg, errorList(f.errs))
	}

	return msg
}

func (f *fault) String() string {
	return f.Error()
}

func (f *fault) Code() Code {
	return f.code
}

func (f *fault) Status() Status {
	return f.status
}

func (f *fault) Message() string {
	return f.message
}

func (f *fault) OrigErr() error {
	switch len(f.errs) {
	case 0:
		return nil
	case 1:
		return f.errs[0]
	default:
		if err, ok := f.errs[0].(Fault); ok {
			return newFault(err.Status(), err.Code(), err.Message(), f.errs[1:])
		}

		return newFault("Errors", "unknown", "multiple errors occurred", f.errs)
	}
}

func newFault(status Status, code Code, message string, errs []error) Fault {
	return &fault{
		status:  status,
		code:    code,
		message: message,
		errs:    errs,
	}
}

func errorList(errs []error) string {
	msg := ""

	if size := len(errs); size > 0 {
		for i := 0; i < size; i++ {
			msg += errs[i].Error()
			if i+1 < size {
				msg += "\n"
			}
		}
	}

	return msg
}
