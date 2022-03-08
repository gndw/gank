package impl

import (
	"context"
	"net/http"
)

func (h *Handler) Get(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
	return h.usecase.Get(ctx)
}
