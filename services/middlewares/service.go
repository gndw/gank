package middlewares

import (
	"net/http"

	"github.com/gndw/gank/model"
)

type Service interface {
	GetHttpMiddleware(f model.Middleware) http.HandlerFunc
	GetRecovererMiddleware(f model.Middleware) model.Middleware
	GetLoggerMiddleware() func(next http.Handler) http.Handler
	GetAuthMiddleware(isActivateAuth bool, f model.Middleware) model.Middleware
	IsAuthMiddlewareValid() (isValid bool)
}
