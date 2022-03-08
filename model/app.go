package model

import (
	"context"

	"github.com/gndw/gank/services/utils/log"
	"go.uber.org/fx"
)

type App struct {
	Lifecycler Lifecycler
}

func (a *App) Run() (err error) {
	return a.Lifecycler.Run()
}

type BuilderOption func(*App) error

type Lifecycler interface {
	// Executed before builder options
	PreConfig() (err error)
	// Executed after builder options
	PostConfig() (err error)
	// Function to start the lifecycler
	Run() (err error)
	AddProviders(providers ...interface{}) (err error)
	AddInvokers(invokers ...interface{}) (err error)
	OverrideLogger(log log.Service) (err error)
}

type Shutdowner fx.Shutdowner
type Lifecycle interface {
	fx.Lifecycle
	AppendOnError(func(context.Context) error)
	ExecuteOnErrors(ctx context.Context)
}
type Hook struct {
	OnStart func(context.Context) error
	OnStop  func(context.Context) error
}

func NewHook(obj Hook) fx.Hook {
	return fx.Hook{
		OnStart: obj.OnStart,
		OnStop:  obj.OnStop,
	}
}

type MyLifecycle struct {
	Fxl            fx.Lifecycle
	OnInvokeErrors []func(context.Context) error
}

func (l *MyLifecycle) Append(hook fx.Hook) {
	l.Fxl.Append(hook)
}

func (l *MyLifecycle) AppendOnError(f func(context.Context) error) {
	l.OnInvokeErrors = append(l.OnInvokeErrors, f)
}

func (l *MyLifecycle) ExecuteOnErrors(ctx context.Context) {
	for i := len(l.OnInvokeErrors) - 1; i >= 0; i-- {
		l.OnInvokeErrors[i](ctx)
	}
}
