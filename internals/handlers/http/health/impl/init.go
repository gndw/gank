package impl

import (
	commonConstant "github.com/gndw/gank/constant"
	internalConstant "github.com/gndw/gank/internals/constant"
	"github.com/gndw/gank/internals/handlers/http/health"

	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/http/server"

	uHealth "github.com/gndw/gank/internals/usecases/health"
)

type Handler struct {
	usecase uHealth.Usecase
}

func New(server server.Service, usecase uHealth.Usecase) (result health.Handler, err error) {
	h := &Handler{
		usecase: usecase,
	}

	server.AddHttpHandler(
		model.AddHTTPRequest{
			Method:   commonConstant.HTTPMethodGet,
			Endpoint: internalConstant.HTTPEndpointHealth,
			Handler:  h.Get,
		},
	)

	server.AddHttpHandler(
		model.AddHTTPRequest{
			Method:         commonConstant.HTTPMethodGet,
			Endpoint:       internalConstant.HTTPEndpointExampleProtectedContent,
			IsActivateAuth: true,
			Handler:        h.GetProtectedExample,
		},
	)

	return h, nil
}
