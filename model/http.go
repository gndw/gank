package model

import (
	"context"
	"errors"
	"net/http"

	"github.com/gndw/gank/constant"
)

type Middleware func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (interface{}, error)

type MiddlewareOptionKey string
type MiddlewareOption struct {
	Key   MiddlewareOptionKey
	Value interface{}
}

type MiddlewaresWithConfig struct {
	Middlewares []func(m Middleware, options ...MiddlewareOption) Middleware
	Config      []MiddlewareOption
}

type AddHTTPRequest struct {
	Method                constant.HTTPMethod
	Endpoint              string
	Middlewares           []func(m Middleware, options ...MiddlewareOption) Middleware
	MiddlewaresWithConfig MiddlewaresWithConfig
	Handler               Middleware
}

func (m *AddHTTPRequest) Validate() error {
	if len(m.MiddlewaresWithConfig.Middlewares) > 0 && len(m.Middlewares) > 0 {
		return errors.New("cannot use double middleware & middlewareWithConfig, please only use either one")
	}
	return nil
}

func (m *AddHTTPRequest) GetMiddlewares() []func(m Middleware, options ...MiddlewareOption) Middleware {
	if len(m.MiddlewaresWithConfig.Middlewares) > 0 {
		return m.MiddlewaresWithConfig.Middlewares
	}
	if len(m.Middlewares) > 0 {
		return m.Middlewares
	}
	return nil
}

type HTTPResponse struct {
	Data  interface{}   `json:"data,omitempty"`
	Error []interface{} `json:"errors,omitempty"`
}
