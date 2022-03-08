package impl

import (
	"context"
	"net/http"
)

func (h *Handler) GetProtectedExample(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
	return "protected content", nil
}
