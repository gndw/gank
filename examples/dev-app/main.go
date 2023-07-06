package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/gndw/gank/services/middlewares"
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
			func(server server.Service, middlewares middlewares.Service) (err error) {

				// adding OK endpoint
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGET,
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
					Method:   constant.HTTPMethodGET,
					Endpoint: "/my-custom-endpoint/slow",
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
						time.Sleep(time.Millisecond * 100)
						return "OK but slow", nil
					},
				})
				if err != nil {
					return err
				}

				// adding OK endpoint
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodPOST,
					Endpoint: "/my-custom-endpoint/ok",
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {

						type TestRequest struct {
							Username string `json:"username"`
							Password string `json:"password"`
						}

						type TestResponse struct {
							UserID int64  `json:"user_id"`
							Token  string `json:"token"`
						}

						var req TestRequest

						exist, rb := contextg.GetRequestBody(ctx)
						if exist {
							err = json.Unmarshal(rb, &req)
							if err != nil {
								return data, err
							}
						}

						return TestResponse{UserID: 69, Token: "private-token"}, nil
					},
				})
				if err != nil {
					return err
				}

				// adding Bad Request endpoint
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGET,
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
					Method:   constant.HTTPMethodGET,
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
					Method:   constant.HTTPMethodGET,
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
					Method:   constant.HTTPMethodGET,
					Endpoint: "/my-custom-endpoint/tracer",
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {

						ctx, tracer := contextg.WithTracer(ctx, contextg.FromFunction(main))
						defer tracer.Finish()
						time.Sleep(time.Millisecond * 100)

						ctx, tracer2 := contextg.WithTracer(ctx, "testing2")
						time.Sleep(time.Millisecond * 100)

						ctx, tracer3 := contextg.WithTracer(ctx, contextg.FromFunction(errorsg.WithOptions))
						time.Sleep(time.Millisecond * 100)

						tracer3.Finish()
						time.Sleep(time.Millisecond * 50)
						tracer2.Finish()
						time.Sleep(time.Millisecond * 50)
						contextg.WithUserID(ctx, 69)

						return "OK", nil
					},
				})
				if err != nil {
					return err
				}

				// adding OK endpoint with custom middlewares
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGET,
					Endpoint: "/my-custom-endpoint/mid-ok",
					Middlewares: middlewares.GetDefaultWith(
						func(m model.Middleware, options ...model.MiddlewareOption) model.Middleware {
							return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
								fmt.Println("hello 1")
								return m(ctx, rw, r)
							}
						},
						func(m model.Middleware, options ...model.MiddlewareOption) model.Middleware {
							return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
								fmt.Println("hello 2")
								return m(ctx, rw, r)
							}
						},
					),
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
						return "OK", nil
					},
				})
				if err != nil {
					return err
				}

				// adding OK endpoint with custom middlewares with config
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGET,
					Endpoint: "/my-custom-endpoint/mid-ok-config",
					MiddlewaresWithConfig: model.MiddlewaresWithConfig{
						// default middleware will enable logging
						Middlewares: middlewares.GetDefault(),
						// create custom config middleware to masking logging sensitive fields in response
						Config: []model.MiddlewareOption{
							middlewares.WithLoggerMiddlewareOption_AdditionalSensitiveFields(
								[]string{"sensitive_message"},
							),
						},
					},
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
						type Response struct {
							Message          string `json:"message"`
							SensitiveMessage string `json:"sensitive_message"`
						}
						return Response{Message: "msg", SensitiveMessage: "s-msg"}, nil
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
