package gank

import (
	"reflect"

	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/db"
	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/http/router"
	"github.com/gndw/gank/services/http/server"
	"github.com/gndw/gank/services/lifecycler"
	"github.com/gndw/gank/services/middlewares"
	"github.com/gndw/gank/services/secret"
	"github.com/gndw/gank/services/utils/hash"
	"github.com/gndw/gank/services/utils/log"
	"github.com/gndw/gank/services/utils/token"

	lifecyclerService "github.com/gndw/gank/services/lifecycler/impl"

	configService "github.com/gndw/gank/services/config/impl"
	dbService "github.com/gndw/gank/services/db/impl"
	envService "github.com/gndw/gank/services/env/impl"
	serverService "github.com/gndw/gank/services/http/server/impl"
	middlewareService "github.com/gndw/gank/services/middlewares/impl"
	secretService "github.com/gndw/gank/services/secret/impl"

	routerService "github.com/gndw/gank/services/http/router/impl"
	hashService "github.com/gndw/gank/services/utils/hash/impl"
	logService "github.com/gndw/gank/services/utils/log/impl"
	tokenService "github.com/gndw/gank/services/utils/token/impl"

	healthHandler "github.com/gndw/gank/internals/handlers/http/health/impl"
	healthUsecase "github.com/gndw/gank/internals/usecases/health/impl"
)

var (
	ConfigKey      reflect.Type = reflect.TypeOf(func(config.Service) {}).In(0)
	EnvKey         reflect.Type = reflect.TypeOf(func(env.Service) {}).In(0)
	ServerKey      reflect.Type = reflect.TypeOf(func(server.Service) {}).In(0)
	MiddlewaresKey reflect.Type = reflect.TypeOf(func(middlewares.Service) {}).In(0)
	SecretKey      reflect.Type = reflect.TypeOf(func(secret.Service) {}).In(0)

	DbKey     reflect.Type = reflect.TypeOf(func(db.Service) {}).In(0)
	RouterKey reflect.Type = reflect.TypeOf(func(router.Service) {}).In(0)
	HashKey   reflect.Type = reflect.TypeOf(func(hash.Service) {}).In(0)
	LogKey    reflect.Type = reflect.TypeOf(func(log.Service) {}).In(0)
	TokenKey  reflect.Type = reflect.TypeOf(func(token.Service) {}).In(0)
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

func GetDefaultInternalProviders() (providers map[reflect.Type]interface{}) {
	return map[reflect.Type]interface{}{
		ConfigKey:      configService.New,
		EnvKey:         envService.New,
		ServerKey:      serverService.New,
		MiddlewaresKey: middlewareService.New,
		SecretKey:      secretService.New,

		RouterKey: routerService.NewGochi,
		HashKey:   hashService.NewBcript,
		LogKey:    logService.NewLogrus,
		TokenKey:  tokenService.NewJWT,
	}
}
func GetDefaultExternalProviders() (providers map[reflect.Type]interface{}) {
	return map[reflect.Type]interface{}{
		DbKey: dbService.NewSqlx,
	}
}

func WithDefaultProviders(options ...ProvidersBuilderOption) (opt model.BuilderOption) {
	return func(app *model.App) (err error) {
		internals := GetDefaultInternalProviders()
		externals := GetDefaultExternalProviders()
		f := func() (providers map[reflect.Type]interface{}) {
			res := make(map[reflect.Type]interface{})
			for k, v := range internals {
				res[k] = v
			}
			for k, v := range externals {
				res[k] = v
			}
			return res
		}

		bi := withDefaultProviders(f, options...)
		err = bi(app)
		if err != nil {
			return err
		}

		return nil
	}
}

func WithDefaultInternalProviders(options ...ProvidersBuilderOption) (opt model.BuilderOption) {
	return func(app *model.App) (err error) {
		bi := withDefaultProviders(GetDefaultInternalProviders, options...)
		err = bi(app)
		if err != nil {
			return err
		}
		return nil
	}
}

func WithDefaultExternalProviders(options ...ProvidersBuilderOption) (opt model.BuilderOption) {
	return func(app *model.App) (err error) {
		bi := withDefaultProviders(GetDefaultExternalProviders, options...)
		err = bi(app)
		if err != nil {
			return err
		}
		return nil
	}
}

func withDefaultProviders(getProviders func() (providers map[reflect.Type]interface{}), options ...ProvidersBuilderOption) (opt model.BuilderOption) {
	return func(app *model.App) (err error) {
		defaultProviders := getProviders()
		for _, option := range options {
			err = option(&defaultProviders)
			if err != nil {
				return err
			}
		}
		for _, p := range defaultProviders {
			err = app.Lifecycler.AddProviders(p)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

type ProvidersBuilderOption func(*map[reflect.Type]interface{}) error

func DisableDefaultProvider(key reflect.Type) (opt ProvidersBuilderOption) {
	return func(providers *map[reflect.Type]interface{}) (err error) {
		_, exist := (*providers)[key]
		if exist {
			delete((*providers), key)
		}
		return nil
	}
}

func OverrideDefaultProvider(key reflect.Type, provider interface{}) (opt ProvidersBuilderOption) {
	return func(providers *map[reflect.Type]interface{}) (err error) {
		(*providers)[key] = provider
		return nil
	}
}

func WithHealthHandler() (opt model.BuilderOption) {
	return func(app *model.App) (err error) {
		err = app.Lifecycler.AddProviders(healthUsecase.New)
		if err != nil {
			return err
		}
		err = app.Lifecycler.AddInvokers(healthHandler.New)
		if err != nil {
			return err
		}
		return nil
	}
}
