package impl

import (
	"strings"

	"github.com/gndw/gank/services/config"
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

func NewGochi(middlewareService middlewares.Service, config config.Service) (router.Service, error) {

	ins := &Service{
		router:            chi.NewRouter(),
		middlewareService: middlewareService,
	}

	ins.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(config.Server.AllowedOrigins, ","),
		AllowedMethods:   strings.Split(config.Server.AllowedMethods, ","),
		AllowedHeaders:   strings.Split(config.Server.AllowedHeaders, ","),
		ExposedHeaders:   strings.Split(config.Server.ExposedHeaders, ","),
		AllowCredentials: config.Server.AllowCredentials,
		MaxAge:           config.Server.CacheMaxAge,
	}))

	ins.router.Use(middleware.Heartbeat("/ping"))

	return ins, nil
}
