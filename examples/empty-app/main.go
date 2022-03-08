package main

import (
	"log"

	"github.com/gndw/gank"
)

func main() {

	err := gank.CreateAndRunApp(gank.DefaultLifecycler(), gank.WithDefaultProviders())
	if err != nil {
		log.Fatal(err)
	}

}
