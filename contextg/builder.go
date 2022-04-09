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
	if existingTracerInterface != nil {
		existingTracer, ok := existingTracerInterface.(*ContextGTracer)
		if ok {
			if existingTracer != nil && existingTracer.FunctionName != "" {
				existingTracer.Child = &newTracer
				return parent, &newTracer
			} else {
				*existingTracer = newTracer
				return parent, existingTracer
			}
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
