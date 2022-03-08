package impl

import (
	"github.com/gndw/gank/internals/usecases/health"
	"github.com/gndw/gank/services/db"
	"go.uber.org/dig"
)

type Usecase struct {
	db db.Service
}

func New(params Parameters) (result health.Usecase, err error) {
	h := &Usecase{
		db: params.Db,
	}
	return h, nil
}

type Parameters struct {
	dig.In
	Db db.Service `optional:"true"`
}
