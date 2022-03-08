package gank

import (
	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/lifecycler"
	lifecyclerService "github.com/gndw/gank/services/lifecycler/impl"
	logService "github.com/gndw/gank/services/utils/log/impl"
)

var (
	CreateApp     func(lc lifecycler.Service, options ...model.BuilderOption) (app *model.App, err error) = functions.CreateApp
	WithProviders func(f ...interface{}) (opt model.BuilderOption)                                        = functions.WithProviders
	WithInvokers  func(f ...interface{}) (opt model.BuilderOption)                                        = functions.WithInvokers
)

func CreateAndRunApp(lc lifecycler.Service, options ...model.BuilderOption) error {

	app, err := CreateApp(lc, options...)
	if err != nil {
		return err
	}

	err = app.Run()
	if err != nil {
		return err
	}

	return nil
}

func DefaultLifecycler() lifecycler.Service {
	lc, _ := lifecyclerService.NewFX(
		lifecyclerService.WithOverrideLogger(logService.NewLogrus),
	)
	return lc
}
