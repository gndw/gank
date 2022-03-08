package impl

import (
	"time"

	"github.com/gndw/gank/services/lifecycler"
	"github.com/gndw/gank/services/utils/log"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

type Service struct {
	container *dig.Container
	app       *fx.App
	fxOptions []fx.Option
	startTime time.Time
	log       log.Service
}

func NewFX(options ...BuilderOption) (lifecycler.Service, error) {
	ins := &Service{}
	for _, opt := range options {
		err := opt(ins)
		if err != nil {
			return nil, err
		}
	}
	return ins, nil
}

type BuilderOption func(*Service) error

func WithOverrideLogger(logGenerator func() (log.Service, error)) (opt BuilderOption) {
	return func(service *Service) (err error) {
		log, err := logGenerator()
		if err != nil {
			return err
		}
		service.log = log
		service.fxOptions = append(service.fxOptions, fx.Logger(log))
		return nil
	}
}
