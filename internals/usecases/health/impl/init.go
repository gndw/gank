package impl

import (
	"github.com/gndw/gank/internals/usecases/health"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/db"
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
	model.In
	Db db.Service `optional:"true"`
}
