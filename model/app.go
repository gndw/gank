package model

import (
	"context"

	"github.com/gndw/gank/services/lifecycler"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

type App struct {
	Lifecycler lifecycler.Service
}

func (a *App) Run() (err error) {
	return a.Lifecycler.Run()
}

type BuilderOption func(*App) error

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

type In struct {
	dig.In
}
