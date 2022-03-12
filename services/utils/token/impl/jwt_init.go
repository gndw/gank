package impl

import (
	"github.com/gndw/gank/services/secret"
	"github.com/gndw/gank/services/utils/log"
	"github.com/gndw/gank/services/utils/token"
)

type Service struct {
	secret string
}

func NewJWT(secret secret.Service, log log.Service) (token.Service, error) {

	key := secret.Token.Key
	if key == "" {
		log.Debugln("token.service> jwt token key not found")
	}

	return &Service{
		secret: key,
	}, nil
}
