package impl

import (
	"context"
	"net/http"

	"github.com/gndw/gank/contextg"
	"github.com/gndw/gank/model"
	"github.com/google/uuid"
)

func (s *Service) GetRequestIDMiddleware(f model.Middleware, options ...model.MiddlewareOption) model.Middleware {
	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		ctx = contextg.WithRequestID(ctx, requestID)

		return f(ctx, rw, r)
	}
}
