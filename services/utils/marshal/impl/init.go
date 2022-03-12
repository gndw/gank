package impl

import (
	"github.com/gndw/gank/services/utils/marshal"
)

type Service struct {
}

func New() (marshal.Service, error) {
	ins := &Service{}
	return ins, nil
}
