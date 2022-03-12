package router

import (
	"net/http"

	"github.com/gndw/gank/model"
)

type Service interface {
	AddHttpHandler(req model.AddHTTPRequest) (err error)
	GetHandler() (handler http.Handler, err error)
	IsAuthRouterValid() (isValid bool)
}
