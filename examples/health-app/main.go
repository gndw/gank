package main

import (
	"log"

	"github.com/gndw/gank"
)

func main() {

	err := gank.CreateAndRunApp(
		gank.DefaultLifecycler(),
		gank.WithDefaultInternalProviders(),
		gank.WithHealthHandler(),
	)
	if err != nil {
		log.Fatal(err)
	}

}
