package contextg

import (
	"context"
	"strings"
	"time"

	"github.com/gndw/gank/constant"
)

var ContextKeyTracer constant.ContextKey = "tracer"

type CustomContext struct {
	ctx    context.Context
	Tracer *ContextGTracer `json:"tracer,omitempty"`
	UserID *int64          `json:"user_id,omitempty"`
}

func (c *CustomContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *CustomContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *CustomContext) Err() error {
	return c.ctx.Err()
}

func (c *CustomContext) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func CreateCustomContext(parent context.Context) (ctx context.Context) {
	customContext := &CustomContext{
		ctx: parent,
	}
	return customContext
}

func GetMetadata(ctx context.Context) (metadata map[string]interface{}) {

	metadata = make(map[string]interface{})

	exist, fs := GetTracerFunctions(ctx)
	if exist {
		metadata["ctx.function_traces"] = strings.Join(fs, "-")
	}

	exist, userID := GetUserID(ctx)
	if exist {
		metadata["ctx.user_id"] = userID
	}

	return metadata
}
