package impl

import (
	"github.com/gndw/gank/services/utils/machinevar"
)

type Service struct {
}

func New() (machinevar.Service, error) {
	ins := &Service{}
	return ins, nil
}
