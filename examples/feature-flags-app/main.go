package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gndw/gank"
	"github.com/gndw/gank/constant"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/featureflags"
	"github.com/gndw/gank/services/http/server"
	"github.com/gndw/gank/services/middlewares"

	featureFlagsService "github.com/gndw/gank/services/featureflags/impl"
)

func main() {

	err := gank.CreateAndRunApp(
		gank.DefaultLifecycler(),
		gank.WithDefaultInternalProviders(),

		// custom providers
		gank.WithProviders(
			// read config file
			config.CreateDevelopmentPreference("github.com", "gndw", "gank"),

			// using feature flag service
			featureFlagsService.NewRemoteJsonURL,
		),

		// add new http endpoint
		gank.WithInvokers(
			func(server server.Service, middlewares middlewares.Service, featureFlag featureflags.Service) (err error) {

				err = server.AddHttpHandler(model.AddHTTPRequest{
					Method:   constant.HTTPMethodGET,
					Endpoint: "/ff",
					Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {

						// reading feature flag
						return featureFlag.GetKey("Example")
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
