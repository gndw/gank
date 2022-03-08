package impl

import (
	"net/http"

	"github.com/gndw/gank/model"
)

func (s *Service) AddHttpHandler(req model.AddHTTPRequest) (err error) {
	s.router.MethodFunc(req.Method, req.Endpoint,
		s.middlewareService.GetHttpMiddleware(
			s.middlewareService.GetAuthMiddleware(req.IsActivateAuth, req.IsBypassIfAuthErrorButReturnStatusUnauthorized,
				req.Handler,
			),
		),
	)
	return nil
}

func (s *Service) GetHandler() (handler http.Handler, err error) {
	return s.router, nil
}
