package impl

import (
	"github.com/gndw/gank/services/http/router"
	"github.com/gndw/gank/services/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	ins.router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// ins.router.Use(middleware.RequestID)
	// ins.router.Use(middlewareService.GetLoggerMiddleware())
	// ins.router.Use(middleware.Recoverer)
	ins.router.Use(middleware.Heartbeat("/ping"))
	// ins.router.Use(render.SetContentType(render.ContentTypeJSON))

	return ins, nil
}
