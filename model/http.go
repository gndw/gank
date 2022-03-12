package model

import (
	"context"
	"net/http"
)

type Middleware func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (interface{}, error)

type AddHTTPRequest struct {
	Method         string
	Endpoint       string
	IsActivateAuth bool
	Handler        Middleware
}

type HTTPResponse struct {
	Data  interface{}   `json:"data,omitempty"`
	Error []interface{} `json:"errors,omitempty"`
}
