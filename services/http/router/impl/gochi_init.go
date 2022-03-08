package impl

import (
	"github.com/gndw/gank/services/http/router"
	"github.com/gndw/gank/services/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Service struct {
	router            *chi.Mux
	middlewareService middlewares.Service
}

func NewGochi(middlewareService middlewares.Service) (router.Service, error) {

	ins := &Service{
		router:            chi.NewRouter(),
		middlewareService: middlewareService,
	}

	ins.router.Use(middlewareService.GetLoggerMiddleware())
	ins.router.Use(middleware.Recoverer)
	ins.router.Use(middleware.Heartbeat("/ping"))
	ins.router.Use(render.SetContentType(render.ContentTypeJSON))

	return ins, nil
}
