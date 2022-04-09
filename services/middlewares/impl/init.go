package impl

import (
	"net/http"
	"time"

	"github.com/gndw/gank/contextg"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/middlewares"
	"github.com/gndw/gank/services/utils/log"
	"github.com/gndw/gank/services/utils/token"
)

type Service struct {
	logService    log.Service
	tokenService  token.Service
	configService config.Service
}

func New(log log.Service, token token.Service, config config.Service) (middlewares.Service, error) {

	if !token.IsValid() {
		log.Debugf("middleware.service> middleware service is receiving invalid token service")
	}

	ins := &Service{
		logService:    log,
		tokenService:  token,
		configService: config,
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
