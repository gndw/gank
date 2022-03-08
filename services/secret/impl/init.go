package impl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gndw/gank/constant"
	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/secret"
)

func New(env env.Service) (secret.Service, error) {

	ins := &Service{}

	if env.IsDevelopment() {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			return nil, fmt.Errorf("GOPATH environment is empty")
		}

		// * make sure your project repository is in GOPATH/src/github.com/gndw/social-story-service

		secretFilePath := path.Join(gopath, "src", "github.com", "gndw", "social-story-service", "files", "var", "secret", "ss.development.secret.json")
		secretByte, err := ioutil.ReadFile(secretFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read secret file in %v with err: %v", secretFilePath, err)
		}
		err = json.Unmarshal(secretByte, ins)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal secret file")
		}
	} else if env.IsStaging() || env.IsProduction() {
		secretString := os.Getenv(constant.ENV_SECRET_NAME)
		err := json.Unmarshal([]byte(secretString), ins)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal secret file")
		}
	}

	return ins, nil
}
