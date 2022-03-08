package health

import (
	"context"
	"net/http"
)

type Handler interface {
	Get(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error)
	GetProtectedExample(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error)
}
