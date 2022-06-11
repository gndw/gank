package impl

import (
	"net/http"

	"github.com/gndw/gank/model"
)

func (s *Service) AddHttpHandler(req model.AddHTTPRequest) (err error) {

	// get default
	if len(req.Middlewares) == 0 {
		req.Middlewares = s.middlewareService.GetDefault()
	}

	// setup middlewares
	handler := req.Handler
	for i := len(req.Middlewares) - 1; i >= 0; i-- {
		handler = req.Middlewares[i](handler)
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
