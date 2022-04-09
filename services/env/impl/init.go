package impl

import (
	"fmt"
	"log"

	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/flag"
	"github.com/gndw/gank/services/utils/machinevar"
)

type Service struct {
	env string
}

func New(params Parameters) (env.Service, error) {

	ins := &Service{}
	allowedEnvs := ins.GetAllowedEnv(params.Preference)
	defer func() {
		log.Printf("starting application with env: %v", ins.Get())
	}()

	// get environment from flag
	env, isValidFromFlag, err := ins.GetEnvNameFromFlag(params.Flag, allowedEnvs)
	if err != nil {
		return nil, err
	} else if isValidFromFlag {
		ins.env = env
		log.Printf("env.service> found env [%v] from flag -env", env)
		return ins, nil
	}

	// get environment from machine variable
	machineKey := ins.GetMachineVarKeyForGetEnv(params.Preference)
	env, isValidFromMachinevar, err := ins.GetEnvNameFromMachinevar(params.Machinevar, machineKey, allowedEnvs)
	if err != nil {
		return nil, err
	} else if isValidFromMachinevar {
		ins.env = env
		log.Printf("env.service> found env [%v] from Machine Environment Variable with key [%v]", env, machineKey)
		return ins, nil
	}

	// use default env
	ins.env = ins.GetDefaultEnv(params.Preference)
	log.Printf("env.service> no env found in flag & machine-var. default env: %v is used", ins.Get())
	return ins, nil
}

func (s *Service) GetAllowedEnv(pref *env.Preference) (allowedEnvs []string) {
	allowedEnvs = env.DEFAULT_ALLOWED_ENV_NAME
	if pref != nil {
		allowedEnvs = append(allowedEnvs, pref.AdditionalEnvs...)
	}
	return allowedEnvs
}

func (s *Service) GetMachineVarKeyForGetEnv(pref *env.Preference) (machinevarKey string) {
	if pref != nil && pref.MachineEnvName != "" {
		return pref.MachineEnvName
	}
	return env.DEFAULT_MACHINE_ENV_NAME
}

func (s *Service) GetDefaultEnv(pref *env.Preference) (machinevarKey string) {
	if pref != nil && pref.DefaultEnv != "" {
		return pref.DefaultEnv
	}
	return env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT
}

func (s *Service) GetEnvNameFromFlag(flag flag.Service, allowedEnvs []string) (env string, isValid bool, err error) {
	if flag.Env != nil && *flag.Env != "" {
		envFromFlag := *flag.Env
		for _, allowedEnv := range allowedEnvs {
			if allowedEnv == envFromFlag {
				return allowedEnv, true, nil
			}
		}
		return env, false, fmt.Errorf("found environment name [%v] from flag -env but this environment is not allowed", envFromFlag)
	}
	return env, false, nil
}

func (s *Service) GetEnvNameFromMachinevar(machinevar machinevar.Service, key string, allowedEnvs []string) (env string, isValid bool, err error) {

	envFromMachine, err := machinevar.GetVar(key)
	if err == nil {
		for _, allowedEnv := range allowedEnvs {
			if allowedEnv == envFromMachine {
				return allowedEnv, true, nil
			}
		}
		return env, false, fmt.Errorf("found environment name [%v] from machine env-var with key [%v] but this environment is not allowed", envFromMachine, key)
	}
	return env, false, nil
}

type Parameters struct {
	model.In
	Flag       flag.Service
	Machinevar machinevar.Service
	Preference *env.Preference `optional:"true"`
}
