package middlewares

import (
	"net/http"

	"github.com/gndw/gank/model"
)

type Service interface {
	GetInitializeMiddleware(f model.Middleware) http.HandlerFunc
	GetLoggerMiddleware(f model.Middleware) model.Middleware
	GetHttpMiddleware(f model.Middleware) model.Middleware
	GetRecovererMiddleware(f model.Middleware) model.Middleware
	GetRequestIDMiddleware(f model.Middleware) model.Middleware
	GetAuthMiddleware(isActivateAuth bool, f model.Middleware) model.Middleware
	IsAuthMiddlewareValid() (isValid bool)
}
