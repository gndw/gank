package impl

import "github.com/gndw/gank/model"

func (s *Service) GetDefault() []func(m model.Middleware, options ...model.MiddlewareOption) model.Middleware {
	return []func(m model.Middleware, options ...model.MiddlewareOption) model.Middleware{
		s.GetLoggerMiddleware,
		s.GetRequestIDMiddleware,
		s.GetHttpMiddleware,
		s.GetRecovererMiddleware,
	}
}

func (s *Service) GetDefaultWith(middlewares ...func(m model.Middleware, options ...model.MiddlewareOption) model.Middleware) []func(m model.Middleware, options ...model.MiddlewareOption) model.Middleware {
	return append(s.GetDefault(), middlewares...)
}
