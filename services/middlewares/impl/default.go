package impl

import "github.com/gndw/gank/model"

func (s *Service) GetDefault() []func(m model.Middleware) model.Middleware {
	return []func(m model.Middleware) model.Middleware{
		s.GetLoggerMiddleware,
		s.GetRequestIDMiddleware,
		s.GetHttpMiddleware,
		s.GetRecovererMiddleware,
	}
}

func (s *Service) GetDefaultWith(middlewares ...func(m model.Middleware) model.Middleware) []func(m model.Middleware) model.Middleware {
	def := s.GetDefault()
	def = append(def, middlewares...)
	return def
}
