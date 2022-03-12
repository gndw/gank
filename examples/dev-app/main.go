package main

import (
	"fmt"
	"log"

	"github.com/gndw/gank"
	"github.com/gndw/gank/services/config"
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

		// test
		gank.WithInvokers(
			func(secretContent secret.Content, configContent config.Content) error {
				fmt.Println("secret", string(secretContent.Value))
				fmt.Println("config", string(configContent.Value))
				return nil
			},
		),
	)
	if err != nil {
		log.Fatal(err)
	}

}
