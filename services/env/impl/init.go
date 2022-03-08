package impl

import (
	"flag"
	"fmt"
	"os"

	"github.com/gndw/gank/constant"
	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/utils/log"
)

type Service struct {
	env string
}

func New(log log.Service) (env.Service, error) {

	ins := &Service{}
	defer func() {
		log.Infof("starting application with env: %v", ins.Get())
	}()

	// *ENV_DEVELOPMENT = "development"
	// *ENV_STAGING     = "staging"
	// *ENV_PRODUCTION  = "production"

	allowedEnvs := map[string]bool{
		constant.ENV_DEVELOPMENT: true,
		constant.ENV_STAGING:     true,
		constant.ENV_PRODUCTION:  true,
	}

	// * ENV_FLAG_NAME = "env"
	// * ENV_VAR_NAME  = "ss-env"

	// checking env from flag
	flagEnv := flag.String(constant.ENV_FLAG_NAME, "", "process environment")
	flag.Parse()
	if flagEnv != nil && *flagEnv != "" {
		if allow, exist := allowedEnvs[*flagEnv]; exist && allow {
			ins.env = *flagEnv
			return ins, nil
		} else {
			return nil, fmt.Errorf("environment variable [%v] found in flag [%v] is not allowed", *flagEnv, constant.ENV_FLAG_NAME)
		}
	}

	// checking env from machine environment variables
	env := os.Getenv(constant.ENV_VAR_NAME)
	if env != "" {
		if allow, exist := allowedEnvs[env]; exist && allow {
			ins.env = env
			return ins, nil
		} else {
			return nil, fmt.Errorf("environment variable [%v] found in machine variable [%v] not allowed", env, constant.ENV_VAR_NAME)
		}
	}

	// if flag & machine env vars is not found, use default 'development'
	ins.env = constant.ENV_DEVELOPMENT

	return ins, nil
}
