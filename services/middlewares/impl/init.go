package impl

import (
	"net/http"

	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/middlewares"
	"github.com/gndw/gank/services/utils/log"
	"github.com/gndw/gank/services/utils/token"
)

type Service struct {
	logService    log.Service
	tokenService  token.Service
	configService config.Service
	logMiddleware func(next http.Handler) http.Handler
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
