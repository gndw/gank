package model

import (
	"context"
	"net/http"

	"github.com/gndw/gank/constant"
)

type Middleware func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (interface{}, error)

type AddHTTPRequest struct {
	Method      constant.HTTPMethod
	Endpoint    string
	Middlewares []func(m Middleware) Middleware
	Handler     Middleware
}

type HTTPResponse struct {
	Data  interface{}   `json:"data,omitempty"`
	Error []interface{} `json:"errors,omitempty"`
}
