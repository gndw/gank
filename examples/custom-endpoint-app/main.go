package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gndw/gank"
	"github.com/gndw/gank/constant"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/http/server"
)

func main() {

	// main function to create your application
	err := gank.CreateAndRunApp(

		// lifecycler object to handle start, stop, and dependencies injection (must use)
		gank.DefaultLifecycler(),

		// inject all default services such as : logger, http server, router, middlewares, token, etc
		// those services will not be called if not used
		gank.WithDefaultInternalProviders(),

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

				// adding Bad Request endpoint
				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGet,
					Endpoint: "/my-custom-endpoint/bad",
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
						return nil, errors.New("bad request response")
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
