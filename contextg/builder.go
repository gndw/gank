package contextg

import (
	"context"

	"github.com/gndw/gank/constant"
)

func WithTracer(parent context.Context, functionName string) (ctx context.Context, tracer *ContextGTracer) {

	// create new tracer
	newTracer := ContextGTracer{}
	newTracer.Start(functionName)

	// check parent tracer
	customContext, ok := parent.(*CustomContext)
	if ok {
		if customContext.Tracer != nil {
			customContext.Tracer.Child = &newTracer
		} else {
			customContext.Tracer = &newTracer
		}
	} else {
		// create new custom context
		customContext = &CustomContext{
			ctx:    parent,
			Tracer: &newTracer,
		}
	}

	return customContext, &newTracer
}

func GetTracerFunctions(ctx context.Context) (isExist bool, functionTraces []string) {

	// get from our custom context first
	customContext, ok := ctx.(*CustomContext)
	if ok && customContext.Tracer != nil {
		functionTraces = []string{}
		currentTrace := customContext.Tracer
		for currentTrace != nil {
			functionTraces = append(functionTraces, currentTrace.FunctionName)
			currentTrace = currentTrace.Child
		}
		return true, functionTraces
	}

	return false, functionTraces
}

func WithUserID(parent context.Context, userID int64) (ctx context.Context) {

	customContext, ok := parent.(*CustomContext)
	if ok {
		customContext.UserID = &userID
		customContext.ctx = context.WithValue(parent, constant.ContextKeyUserID, userID)
	} else {
		// create new custom context
		customContext = &CustomContext{
			ctx:    context.WithValue(parent, constant.ContextKeyUserID, userID),
			UserID: &userID,
		}
	}

	return customContext
}

func GetUserID(ctx context.Context) (isExist bool, userID int64) {

	// get from our custom context first
	customContext, ok := ctx.(*CustomContext)
	if ok && customContext.UserID != nil {
		return true, *customContext.UserID
	}

	// try to get user ID from default context
	userIDi := ctx.Value(constant.ContextKeyUserID)
	if userIDi != nil {
		userID, ok := userIDi.(int64)
		if ok {
			WithUserID(ctx, userID)
			return true, userID
		}
	}

	return false, 0
}
