package impl

import "github.com/gndw/gank/services/utils/hash"

type Service struct {
}

func NewBcript() (hash.Service, error) {
	ins := &Service{}
	return ins, nil
}
