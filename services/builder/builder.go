package builder

import (
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/utils/log"
)

func CreateApp(lc model.Lifecycler, options ...model.BuilderOption) (app *model.App, err error) {

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

func OverrideLogger(logGenerator func() (log.Service, error)) (opt model.BuilderOption) {
	return func(app *model.App) (err error) {
		log, err := logGenerator()
		if err != nil {
			return err
		}
		return app.Lifecycler.OverrideLogger(log)
	}
}

func SetProvidersAndInvokers(providers []interface{}, invokers []interface{}) (opt model.BuilderOption) {
	return func(app *model.App) (err error) {
		err = app.Lifecycler.AddProviders(providers...)
		if err != nil {
			return err
		}

		err = app.Lifecycler.AddInvokers(invokers...)
		if err != nil {
			return err
		}

		return nil
	}
}
