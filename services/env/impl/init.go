package impl

import (
	"fmt"
	"os"

	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/flag"
	"github.com/gndw/gank/services/utils/log"
)

type Service struct {
	env            string
	defaultEnv     string
	machineEnvName string
	allowedEnvs    map[string]bool
}

func New(params Parameters) (env.Service, error) {

	ins := &Service{}
	ins.PopulateDataFromPreference(params.Preference)

	defer func() {
		params.Log.Infof("starting application with env: %v", ins.Get())
	}()

	isExistFromFlag, err := ins.PopulateEnvNameFromFlag(params.Flag.Env)
	if err != nil {
		return nil, err
	} else if isExistFromFlag {
		params.Log.Debugf("env.service> found env: %v from flag", ins.Get())
		return ins, nil
	}

	isExistFromMachineEnvVar, err := ins.PopulateEnvNameFromEnvMachineVar()
	if err != nil {
		return nil, err
	} else if isExistFromMachineEnvVar {
		params.Log.Debugf("env.service> found env: %v from machine env-var", ins.Get())
		return ins, nil
	}

	// use default env
	ins.env = ins.defaultEnv
	params.Log.Debugf("env.service> default env: %v is used", ins.Get())
	return ins, nil
}

func (s *Service) PopulateDataFromPreference(pref *env.Preference) {

	s.defaultEnv = env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT
	s.machineEnvName = env.DEFAULT_MACHINE_ENV_NAME
	s.allowedEnvs = make(map[string]bool)
	for _, env := range env.DEFAULT_ALLOWED_ENV_NAME {
		s.allowedEnvs[env] = true
	}

	if pref != nil {
		if functions.IsAllNonEmpty(pref.DefaultEnv) {
			s.defaultEnv = pref.DefaultEnv
		}
		if functions.IsAllNonEmpty(pref.MachineEnvName) {
			s.machineEnvName = pref.MachineEnvName
		}
		for _, addEnv := range pref.AdditionalEnvs {
			s.allowedEnvs[addEnv] = true
		}
	}

}

func (s *Service) PopulateEnvNameFromFlag(flagEnv *string) (isValid bool, err error) {
	if flagEnv != nil && *flagEnv != "" {
		if allow, exist := s.allowedEnvs[*flagEnv]; exist && allow {
			s.env = *flagEnv
			return true, nil
		} else {
			return false, fmt.Errorf("environment variable [%v] found is not allowed", *flagEnv)
		}
	}
	return false, nil
}

func (s *Service) PopulateEnvNameFromEnvMachineVar() (isValid bool, err error) {
	env := os.Getenv(s.machineEnvName)
	if env != "" {
		if allow, exist := s.allowedEnvs[env]; exist && allow {
			s.env = env
			return true, nil
		} else {
			return false, fmt.Errorf("environment variable [%v] found in machine variable [%v] not allowed", env, s.machineEnvName)
		}
	}
	return false, nil
}

type Parameters struct {
	model.In
	Log        log.Service
	Flag       flag.Service
	Preference *env.Preference `optional:"true"`
}
