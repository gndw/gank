package functions

import (
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/lifecycler"
)

func CreateApp(lc lifecycler.Service, options ...model.BuilderOption) (app *model.App, err error) {

	app = &model.App{
		Lifecycler: lc,
	}

	err = app.Lifecycler.PreConfig()
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		err = option(app)
		if err != nil {
			return nil, err
		}
	}

	err = app.Lifecycler.PostConfig()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func WithProviders(providers ...interface{}) (opt model.BuilderOption) {
	return func(app *model.App) (err error) {
		err = app.Lifecycler.AddProviders(providers...)
		if err != nil {
			return err
		}
		return nil
	}
}

func WithInvokers(invokers ...interface{}) (opt model.BuilderOption) {
	return func(app *model.App) (err error) {
		err = app.Lifecycler.AddInvokers(invokers...)
		if err != nil {
			return err
		}
		return nil
	}
}
