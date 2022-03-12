package server

import (
	"github.com/gndw/gank/model"
)

type Service interface {
	AddHttpHandler(req model.AddHTTPRequest) (err error)
	IsAuthRouterValid() (isValid bool)
}
