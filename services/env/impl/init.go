package impl

import (
	"fmt"

	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/flag"
	"github.com/gndw/gank/services/utils/log"
	"github.com/gndw/gank/services/utils/machinevar"
)

type Service struct {
	env            string
	isReleaseLevel bool
}

func New(params Parameters) (env.Service, error) {

	ins := &Service{}
	allowedEnvs := ins.GetAllowedEnv(params.Preference)
	defer func() {
		params.Log.Infof("starting application with env: %v", ins.Get())
	}()

	// get environment from flag
	env, isValidFromFlag, err := ins.GetEnvNameFromFlag(params.Flag, allowedEnvs)
	if err != nil {
		return nil, err
	} else if isValidFromFlag {
		ins.env = env.EnvName
		ins.isReleaseLevel = env.IsReleaseLevel
		params.Log.Debugf("env.service> found env [%v] with release level [%v] from flag -env", env.EnvName, env.IsReleaseLevel)
		return ins, nil
	}

	// get environment from machine variable
	machineKey := ins.GetMachineVarKeyForGetEnv(params.Preference)
	env, isValidFromMachinevar, err := ins.GetEnvNameFromMachinevar(params.Machinevar, machineKey, allowedEnvs)
	if err != nil {
		return nil, err
	} else if isValidFromMachinevar {
		ins.env = env.EnvName
		ins.isReleaseLevel = env.IsReleaseLevel
		params.Log.Debugf("env.service> found env [%v] with release level [%v] from Machine Environment Variable with key [%v]", env.EnvName, env.IsReleaseLevel, machineKey)
		return ins, nil
	}

	// use default env
	ins.env = ins.GetDefaultEnv(params.Preference)
	params.Log.Debugf("env.service> no env found in flag & machine-var. default env: %v is used", ins.Get())
	return ins, nil
}

func (s *Service) GetAllowedEnv(pref *env.Preference) (allowedEnvs []env.EnvLevel) {
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

func (s *Service) GetEnvNameFromFlag(flag flag.Service, allowedEnvs []env.EnvLevel) (envLevel env.EnvLevel, isValid bool, err error) {
	if flag.Env != nil && *flag.Env != "" {
		envFromFlag := *flag.Env
		for _, allowedEnv := range allowedEnvs {
			if allowedEnv.EnvName == envFromFlag {
				return allowedEnv, true, nil
			}
		}
		return envLevel, false, fmt.Errorf("found environment name [%v] from flag -env but this environment is not allowed", envFromFlag)
	}
	return envLevel, false, nil
}

func (s *Service) GetEnvNameFromMachinevar(machinevar machinevar.Service, key string, allowedEnvs []env.EnvLevel) (envLevel env.EnvLevel, isValid bool, err error) {

	envFromMachine, err := machinevar.GetVar(key)
	if err == nil {
		for _, allowedEnv := range allowedEnvs {
			if allowedEnv.EnvName == envFromMachine {
				return allowedEnv, true, nil
			}
		}
		return envLevel, false, fmt.Errorf("found environment name [%v] from machine env-var with key [%v] but this environment is not allowed", envFromMachine, key)
	}
	return envLevel, false, nil
}

type Parameters struct {
	model.In
	Log        log.Service
	Flag       flag.Service
	Machinevar machinevar.Service
	Preference *env.Preference `optional:"true"`
}
