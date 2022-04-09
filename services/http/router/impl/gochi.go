package impl

import (
	"errors"
	"net/http"

	"github.com/gndw/gank/model"
)

func (s *Service) AddHttpHandler(req model.AddHTTPRequest) (err error) {

	if req.IsActivateAuth && !s.middlewareService.IsAuthMiddlewareValid() {
		return errors.New("cannot add router with activate auth when auth middleware is invalid")
	}

	s.router.MethodFunc(req.Method, req.Endpoint,
		s.middlewareService.GetInitializeMiddleware(
			s.middlewareService.GetLoggerMiddleware(
				s.middlewareService.GetRequestIDMiddleware(
					s.middlewareService.GetHttpMiddleware(
						s.middlewareService.GetRecovererMiddleware(
							s.middlewareService.GetAuthMiddleware(req.IsActivateAuth,
								req.Handler,
							),
						),
					),
				),
			),
		),
	)
	return nil
}

func (s *Service) GetHandler() (handler http.Handler, err error) {
	return s.router, nil
}

func (s *Service) IsAuthRouterValid() (isValid bool) {
	return s.middlewareService.IsAuthMiddlewareValid()
}
