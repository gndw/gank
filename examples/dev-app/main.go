package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gndw/gank"
	"github.com/gndw/gank/constant"
	"github.com/gndw/gank/contextg"
	"github.com/gndw/gank/errorsg"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/http/server"
	"github.com/gndw/gank/services/secret"
)

func main() {

	// main function to create your application
	err := gank.CreateAndRunApp(

		// lifecycler object to handle start, stop, and dependencies injection (must use)
		gank.DefaultLifecycler(),

		// inject all default services such as : logger, http server, router, middlewares, token, etc
		// those services will not be called if not used
		gank.WithDefaultInternalProviders(),

		// use this option to start http health service
		gank.WithHealthHandler(),

		// create custom default env
		// gank.WithProviders(func() (*env.Preference, error) {
		// 	return &env.Preference{
		// 		DefaultEnv: "yahoo",
		// 	}, nil
		// }),

		// add custom config and secret file
		gank.WithProviders(
			config.CreateDevelopmentPreference("github.com", "gndw", "gank"),
			secret.CreateDevelopmentPreference("github.com", "gndw", "gank"),
		),

		// // test
		// gank.WithInvokers(
		// 	func(secretContent secret.Content, configContent config.Content) error {
		// 		fmt.Println("secret", string(secretContent.Value))
		// 		fmt.Println("config", string(configContent.Value))
		// 		return nil
		// 	},
		// ),

		// add new http endpoint
		gank.WithInvokers(
			func(server server.Service) (err error) {

				// adding OK endpoint
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGet,
					Endpoint: "/my-custom-endpoint/ok",
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
						return "OK", nil
					},
				})
				if err != nil {
					return err
				}

				// adding OK endpoint but slow
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGet,
					Endpoint: "/my-custom-endpoint/slow",
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
						time.Sleep(time.Millisecond * 100)
						return "OK but slow", nil
					},
				})
				if err != nil {
					return err
				}

				// adding Bad Request endpoint
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGet,
					Endpoint: "/my-custom-endpoint/bad",
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
						return nil, errorsg.BadRequestWithOptions(errors.New("bad request response"), errorsg.WithPrivateIdentifier("pipipi"))
					},
				})
				if err != nil {
					return err
				}

				// adding Internal Server Error endpoint
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGet,
					Endpoint: "/my-custom-endpoint/error",
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
						return nil, errorsg.BadRequestWithOptions(errors.New("internal error"), errorsg.WithType(errorsg.ErrorTypeInternalServerError))
					},
				})
				if err != nil {
					return err
				}

				// adding Panic endpoint
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGet,
					Endpoint: "/my-custom-endpoint/panic",
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
						var test interface{}
						number := test.(int64)
						return number, nil
					},
				})
				if err != nil {
					return err
				}

				// adding endpoint with tracer
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGet,
					Endpoint: "/my-custom-endpoint/tracer",
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {

						ctx, tracer := contextg.WithTracer(ctx, contextg.FromFunction(main))
						defer tracer.Finish()
						time.Sleep(time.Millisecond * 100)

						ctx, tracer2 := contextg.WithTracer(ctx, "testing2")
						time.Sleep(time.Millisecond * 100)
						tracer2.Finish()

						time.Sleep(time.Millisecond * 50)

						contextg.WithUserID(ctx, 69)

						return "OK", nil
					},
				})
				if err != nil {
					return err
				}

				return nil
			},
		),
	)
	if err != nil {
		log.Fatal(err)
	}

}
