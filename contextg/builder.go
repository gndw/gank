package contextg

import (
	"context"
	"time"

	"github.com/gndw/gank/constant"
)

func WithTracer(parent context.Context, functionName string) (ctx context.Context, tracer *ContextGTracer) {

	newTracer := ContextGTracer{}
	newTracer.Start(functionName)

	existingTracerInterface := parent.Value(ContextKeyTracer)
	for existingTracerInterface != nil {
		existingTracer, ok := existingTracerInterface.(*ContextGTracer)
		if ok {
			if existingTracer != nil && existingTracer.FunctionName != "" {
				if existingTracer.Child != nil {
					existingTracerInterface = existingTracer.Child
				} else {
					existingTracer.Child = &newTracer
					return parent, &newTracer
				}
			} else {
				*existingTracer = newTracer
				return parent, existingTracer
			}
		} else {
			break
		}

	}

	return parent, nil
}

func GetTracerFunctions(ctx context.Context) (isExist bool, functionTraces []string) {

	existingTracerInterface := ctx.Value(ContextKeyTracer)
	if existingTracerInterface != nil {
		existingTracer, ok := existingTracerInterface.(*ContextGTracer)
		if ok {
			if existingTracer != nil && existingTracer.FunctionName != "" {
				functionTraces = []string{}
				currentTrace := existingTracer
				for currentTrace != nil {
					functionTraces = append(functionTraces, currentTrace.FunctionName)
					currentTrace = currentTrace.Child
				}
				return true, functionTraces
			}
		}
	}

	return false, nil
}

func WithUserID(parent context.Context, userID int64) (ctx context.Context) {

	obj := parent.Value(constant.ContextKeyUserID)
	if obj != nil {
		pointer, ok := obj.(*int64)
		if ok {
			*pointer = userID
		}
	}

	return parent
}

func GetUserID(ctx context.Context) (isExist bool, userID int64) {

	obj := ctx.Value(constant.ContextKeyUserID)
	if obj != nil {
		pointer, ok := obj.(*int64)
		if ok && *pointer != 0 {
			return true, *pointer
		}
	}

	return false, 0
}

func WithRequestID(parent context.Context, requestID string) (ctx context.Context) {

	obj := parent.Value(constant.ContextKeyRequestID)
	if obj != nil {
		pointer, ok := obj.(*string)
		if ok {
			*pointer = requestID
		}
	}

	return parent
}

func GetRequestID(ctx context.Context) (isExist bool, requestID string) {

	obj := ctx.Value(constant.ContextKeyRequestID)
	if obj != nil {
		pointer, ok := obj.(*string)
		if ok && *pointer != "" {
			return true, *pointer
		}
	}

	return false, ""
}

func WithIncomingTime(parent context.Context, t time.Time) (ctx context.Context) {

	obj := parent.Value(constant.ContextKeyIncomingTime)
	if obj != nil {
		pointer, ok := obj.(*time.Time)
		if ok {
			*pointer = t
		}
	}

	return parent
}

func GetIncomingTime(ctx context.Context) (isExist bool, t time.Time) {

	obj := ctx.Value(constant.ContextKeyIncomingTime)
	if obj != nil {
		pointer, ok := obj.(*time.Time)
		if ok && !(*pointer).IsZero() {
			return true, *pointer
		}
	}

	return false, t
}

func WithRequestBody(parent context.Context, body []byte) (ctx context.Context) {

	obj := parent.Value(constant.ContextKeyRequestBody)
	if obj != nil {
		pointer, ok := obj.(*[]byte)
		if ok {
			*pointer = body
		}
	}

	return parent
}

func GetRequestBody(ctx context.Context) (isExist bool, body []byte) {

	obj := ctx.Value(constant.ContextKeyRequestBody)
	if obj != nil {
		pointer, ok := obj.(*[]byte)
		if ok && len(*pointer) > 0 {
			return true, *pointer
		}
	}

	return false, nil
}

func WithCustomData(parent context.Context, key string, value interface{}) (ctx context.Context) {

	obj := parent.Value(constant.ContextKeyCustomData)
	if obj != nil {
		pointer, ok := obj.(*map[string]interface{})
		if ok {
			if (*pointer) == nil {
				(*pointer) = make(map[string]interface{})
			}
			(*pointer)[key] = value
		}
	}

	return parent
}

func GetCustomData(ctx context.Context) (isExist bool, data map[string]interface{}) {

	obj := ctx.Value(constant.ContextKeyCustomData)
	if obj != nil {
		pointer, ok := obj.(*map[string]interface{})
		if ok && *pointer != nil {
			return true, *pointer
		}
	}

	return false, nil
}
