package impl

import (
	"net/http"
	"time"

	"github.com/gndw/gank/contextg"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/middlewares"
	"github.com/gndw/gank/services/utils/log"
	"github.com/gndw/gank/services/utils/token"
)

type Service struct {
	logService    log.Service
	configService config.Service
	envService    env.Service
}

func New(log log.Service, token token.Service, config config.Service, env env.Service) (middlewares.Service, error) {

	ins := &Service{
		logService:    log,
		configService: config,
		envService:    env,
	}
	return ins, nil
}

func (s *Service) GetInitializeMiddleware(f model.Middleware) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		// initialize context
		ctx := contextg.CreateCustomContext(r.Context())

		// record incoming time
		ctx = contextg.WithIncomingTime(ctx, time.Now())

		f(ctx, rw, r)
	}
}
