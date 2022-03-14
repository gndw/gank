package impl

import (
	"github.com/gndw/gank/services/utils/io"
)

type Service struct {
}

func New() (io.Service, error) {
	ins := &Service{}
	return ins, nil
}
