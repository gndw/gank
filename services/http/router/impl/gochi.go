package impl

import (
	"net/http"

	"github.com/gndw/gank/model"
)

func (s *Service) AddHttpHandler(req model.AddHTTPRequest) (err error) {

	// validation
	if err := req.Validate(); err != nil {
		return err
	}
	middlewares := req.GetMiddlewares()

	// get default
	if len(middlewares) == 0 {
		middlewares = s.middlewareService.GetDefault()
	}

	// setup middlewares
	handler := req.Handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler, req.MiddlewaresWithConfig.Config...)
	}

	// adding to router
	s.router.MethodFunc(string(req.Method), req.Endpoint,
		s.middlewareService.GetInitializeMiddleware(
			handler,
		),
	)

	return nil
}

func (s *Service) GetHandler() (handler http.Handler, err error) {
	return s.router, nil
}
