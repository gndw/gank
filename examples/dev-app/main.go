package main

import (
	"log"

	"github.com/gndw/gank"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/env"
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

		// gank.WithProviders(func() (*env.Preference, error) {
		// 	return &env.Preference{
		// 		DefaultEnv: "yahoo",
		// 	}, nil
		// }),

		gank.WithProviders(
			config.CreatePreference(config.Preference{
				EnvFilePaths: map[string][]string{
					env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT: config.GetDefaultFilePathUsingRepositoryPath("github.com", "gndw", "social-story-service"),
				},
			}),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

}
