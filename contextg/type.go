package contextg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gndw/gank/constant"
)

var ContextKeyTracer constant.ContextKey = "tracer"

func CreateCustomContext(ctx context.Context) context.Context {

	// initialize all value to create a pointer
	ctx = context.WithValue(ctx, constant.ContextKeyUserID, new(int64))
	ctx = context.WithValue(ctx, constant.ContextKeyRequestID, new(string))
	ctx = context.WithValue(ctx, ContextKeyTracer, new(ContextGTracer))
	ctx = context.WithValue(ctx, constant.ContextKeyIncomingTime, new(time.Time))
	ctx = context.WithValue(ctx, constant.ContextKeyRequestBody, new([]byte))
	ctx = context.WithValue(ctx, constant.ContextKeyCustomData, new(map[string]interface{}))
	return ctx
}

func GetMetadata(ctx context.Context) (metadata map[string]interface{}) {

	metadata = make(map[string]interface{})

	exist, fs := GetTracerFunctions(ctx)
	if exist {
		metadata["ctx.function_traces"] = strings.Join(fs, "|")
	}

	exist, userID := GetUserID(ctx)
	if exist {
		metadata["ctx.user_id"] = userID
	}

	exist, requestID := GetRequestID(ctx)
	if exist {
		metadata["ctx.request_id"] = requestID
	}

	exist, customData := GetCustomData(ctx)
	if exist {
		for k, v := range customData {
			metadata[fmt.Sprintf("ctx.custom.%v.", k)] = v
		}
	}

	return metadata
}
