package impl

import (
	"flag"
	"fmt"
	"os"

	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/utils/log"
)

type Service struct {
	env            string
	defaultEnv     string
	flagNameEnv    string
	machineEnvName string
	allowedEnvs    map[string]bool
}

func New(params Parameters) (env.Service, error) {

	ins := &Service{}
	ins.PopulateDataFromPreference(params.Preference)

	defer func() {
		params.Log.Infof("starting application with env: %v", ins.Get())
	}()

	isExistFromFlag, err := ins.PopulateEnvNameFromFlag()
	if err != nil {
		return nil, err
	} else if isExistFromFlag {
		return ins, nil
	}

	isExistFromMachineEnvVar, err := ins.PopulateEnvNameFromEnvMachineVar()
	if err != nil {
		return nil, err
	} else if isExistFromMachineEnvVar {
		return ins, nil
	}

	// use default env
	ins.env = ins.defaultEnv
	return ins, nil
}

func (s *Service) PopulateDataFromPreference(pref *env.Preference) {

	s.defaultEnv = env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT
	s.flagNameEnv = env.DEFAULT_FLAG_NAME_ENV
	s.machineEnvName = env.DEFAULT_MACHINE_ENV_NAME
	allowedEnvs := map[string]bool{
		env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT: true,
		env.DEFAULT_ENV_NAME_ENV_STAGING:     true,
		env.DEFAULT_ENV_NAME_ENV_PRODUCTION:  true,
	}

	if pref != nil {
		if functions.IsAllNonEmpty(pref.DefaultEnv) {
			s.defaultEnv = pref.DefaultEnv
		}
		if functions.IsAllNonEmpty(pref.FlagNameEnv) {
			s.flagNameEnv = pref.FlagNameEnv
		}
		if functions.IsAllNonEmpty(pref.MachineEnvName) {
			s.machineEnvName = pref.MachineEnvName
		}
		for _, addEnv := range pref.AdditionalEnvs {
			allowedEnvs[addEnv] = true
		}
	}

}

func (s *Service) PopulateEnvNameFromFlag() (isValid bool, err error) {
	flagEnv := flag.String(s.flagNameEnv, "", "process environment")
	flag.Parse()
	if flagEnv != nil && *flagEnv != "" {
		if allow, exist := s.allowedEnvs[*flagEnv]; exist && allow {
			s.env = *flagEnv
			return true, nil
		} else {
			return false, fmt.Errorf("environment variable [%v] found in flag [%v] is not allowed", *flagEnv, s.flagNameEnv)
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
	Preference *env.Preference `optional:"true"`
}
