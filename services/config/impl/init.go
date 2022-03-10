package impl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/env"
)

func New(env env.Service) (config.Service, error) {

	ins := &Service{}
	configFilePath := ""

	if env.IsDevelopment() {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			return nil, fmt.Errorf("GOPATH environment is empty")
		}

		// * make sure your project repository is in GOPATH/src/github.com/gndw/social-story-service

		configFilePath = path.Join(gopath, "src", "github.com", "gndw", "social-story-service", "files", "var", "config", "ss.development.config.yaml")
	} else if env.IsStaging() || env.IsProduction() {

		configFileName := ""
		if env.IsStaging() {
			configFileName = "ss.staging.config.yaml"
		} else if env.IsProduction() {
			configFileName = "ss.production.config.yaml"
		}

		// this is specific for HEROKU server
		serverFolderPath := "/app"
		configFilePath = path.Join(serverFolderPath, "files", "var", "config", configFileName)
	}

	configByte, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file in %v with err: %v", configFilePath, err)
	}
	err = yaml.Unmarshal(configByte, ins)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file")
	}
	return ins, nil
}
