package impl

import (
	"time"

	"github.com/gndw/gank/model"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

type Service struct {
	container *dig.Container
	app       *fx.App
	fxOptions []fx.Option
	startTime time.Time
}

func NewFX() model.Lifecycler {
	return &Service{}
}
