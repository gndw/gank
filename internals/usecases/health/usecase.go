package health

import (
	"context"

	"github.com/gndw/gank/internals/model"
)

type Usecase interface {
	Get(ctx context.Context) (result model.Health, err error)
}
