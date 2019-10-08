package http

import (
	"context"

	"go.uber.org/zap"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/labstack/echo"
	"github.com/samygp/edgex-health-alerts/log"
)

type (
	contextKey string

	// Store defines a generic map.
	Store map[string]interface{}
)

const (
	loggerReqID      = "req_id"
	contextKeyID     = "ctx_id"
	contextKeyEcho   = "ctx_echo"
	contextKeyLambda = "ctx_lambda"
	contextKeyLogger = "ctx_logger"
)

var (
	ctxKey     = contextKey(contextKeyID)
	nullLogger = zap.NewNop().Sugar()
)

// // Get store
// func (v Store) Get(key string) interface{} {
// 	return v[key]
// }

// NewContextFromLambda returns a new instance of Context.
func NewContextFromLambda(base context.Context, ctx *lambdacontext.LambdaContext) context.Context {
	id := ctx.AwsRequestID
	store := make(Store)

	store[contextKeyID] = id
	store[contextKeyLambda] = ctx
	store[contextKeyLogger] = log.Logger.With(zap.String(loggerReqID, id))

	if base == nil {
		base = context.Background()
	}

	return context.WithValue(base, ctxKey, store)
}

// NewContextFromEcho returns a new instance of Context.
func NewContextFromEcho(base context.Context, ctx echo.Context) context.Context {
	id := ctx.Response().Header().Get(echo.HeaderXRequestID)
	store := make(Store)

	store[contextKeyID] = id
	store[contextKeyEcho] = ctx
	store[contextKeyLogger] = log.Logger.With(zap.String(loggerReqID, id))

	if base == nil {
		base = context.Background()
	}

	return context.WithValue(base, ctxKey, store)
}

// ContextLogger returns Logger out of context.
func ContextLogger(ctx context.Context) *zap.SugaredLogger {
	if store := fromContext(ctx); store != nil {
		l, ok := store[contextKeyLogger].(*zap.SugaredLogger)
		if ok && l != nil {
			return l
		}
	}

	return nullLogger
}

// ContextID returns ID ouf of context.
func ContextID(ctx context.Context) string {
	if store := fromContext(ctx); store != nil {
		if id := store[contextKeyID]; id != nil {
			return id.(string)
		}
	}

	return ""
}

// ContextRealIP returns IP address ouf of context.
func ContextRealIP(ctx context.Context) string {
	if store := fromContext(ctx); store != nil {
		ctx, ok := store[contextKeyEcho].(echo.Context)
		if ok && ctx != nil {
			return ctx.RealIP()
		}
	}

	return ""
}

func fromContext(ctx context.Context) Store {
	store, ok := ctx.Value(ctxKey).(Store)
	if !ok || store == nil {
		return nil
	}

	return store
}
