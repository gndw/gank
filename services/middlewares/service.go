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
	GetDefault() []func(m model.Middleware) model.Middleware
	GetDefaultWith(middlewares ...func(m model.Middleware) model.Middleware) []func(m model.Middleware) model.Middleware
}

type Auth interface {
	GetAuthMiddleware(f model.Middleware) model.Middleware
}
