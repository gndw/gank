package impl

import (
	"errors"

	"github.com/gndw/gank/services/secret"
	"github.com/gndw/gank/services/utils/token"
)

type Service struct {
	secret string
}

func NewJWT(secret secret.Service) (token.Service, error) {

	key := secret.GetToken().GetKey()
	if key == "" {
		return nil, errors.New("jwt token key not found")
	}

	ins := &Service{
		secret: key,
	}
	return ins, nil
}
