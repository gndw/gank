package impl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/env"
)

type ConfigInternalData struct {
	filePathDevelopment          string
	mustValidFilePathDevelopment bool

	filePathStaging          string
	mustValidFilePathStaging bool

	filePathProduction          string
	mustValidFilePathProduction bool

	additionalEnvFilePaths map[string]string
}

func New(params Parameters) (config.Service, error) {

	ins := &Service{}

	err := ins.PopulateDataFromPreference(params.Preference)
	if err != nil {
		return nil, err
	}

	if params.Env.IsDevelopment() {
		err = ins.PopulateDataFromConfigFilePath(ins.filePathDevelopment)
		if ins.mustValidFilePathDevelopment && err != nil {
			return nil, err
		}
	} else if params.Env.IsStaging() {
		err = ins.PopulateDataFromConfigFilePath(ins.filePathStaging)
		if ins.mustValidFilePathStaging && err != nil {
			return nil, err
		}
	} else if params.Env.IsProduction() {
		err = ins.PopulateDataFromConfigFilePath(ins.filePathProduction)
		if ins.mustValidFilePathProduction && err != nil {
			return nil, err
		}
	} else {
		currentEnv := params.Env.Get()
		path, exist := ins.additionalEnvFilePaths[currentEnv]
		if exist {
			err = ins.PopulateDataFromConfigFilePath(path)
			if err != nil {
				return nil, err
			}
		}
	}

	return ins, nil
}

func (s *Service) PopulateDataFromPreference(pref *config.Preference) (err error) {

	defaultDevFilePath, _ := s.GetPathFromArray(config.DEFAULT_FILE_PATH_ON_DEVELOPMENT_SERVER)
	s.filePathDevelopment = defaultDevFilePath

	defaultStagFilePath, _ := s.GetPathFromArray(config.DEFAULT_FILE_PATH_ON_STAGING_SERVER)
	s.filePathStaging = defaultStagFilePath

	defaultProdFilePath, _ := s.GetPathFromArray(config.DEFAULT_FILE_PATH_ON_PRODUCTION_SERVER)
	s.filePathProduction = defaultProdFilePath

	if pref != nil {
		if len(pref.FilePathDevelopment) > 0 {
			if functions.IsAllNonEmpty(pref.FilePathDevelopment...) {
				path, err := s.GetPathFromArray(pref.FilePathDevelopment)
				if err != nil {
					return err
				}
				s.filePathDevelopment = path
				s.mustValidFilePathDevelopment = true
			} else {
				return fmt.Errorf("config file path cannot be empty : %v", pref.FilePathDevelopment)
			}
		}
		if len(pref.FilePathStaging) > 0 {
			if functions.IsAllNonEmpty(pref.FilePathStaging...) {
				path, err := s.GetPathFromArray(pref.FilePathStaging)
				if err != nil {
					return err
				}
				s.filePathStaging = path
				s.mustValidFilePathStaging = true
			} else {
				return fmt.Errorf("config file path cannot be empty : %v", pref.FilePathStaging)
			}
		}
		if len(pref.FilePathProduction) > 0 {
			if functions.IsAllNonEmpty(pref.FilePathProduction...) {
				path, err := s.GetPathFromArray(pref.FilePathProduction)
				if err != nil {
					return err
				}
				s.filePathProduction = path
				s.mustValidFilePathProduction = true
			} else {
				return fmt.Errorf("config file path cannot be empty : %v", pref.FilePathProduction)
			}
		}
		for env, addPath := range pref.AdditionalEnvFilePaths {
			if functions.IsAllNonEmpty(addPath...) {
				path, err := s.GetPathFromArray(addPath)
				if err != nil {
					return err
				}
				s.additionalEnvFilePaths[env] = path
			} else {
				return fmt.Errorf("config file path cannot be empty : %v", addPath)
			}
		}
	}
	return nil
}

func (s *Service) PopulateDataFromConfigFilePath(path string) (err error) {
	if path == "" {
		return errors.New("config file path cannot be empty")
	}
	configByte, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file in %v with err: %v", path, err)
	}
	err = yaml.Unmarshal(configByte, s)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config file")
	}
	return nil
}

func (s *Service) GetPathFromArray(pathArray []string) (string, error) {

	sanitizedArray := []string{}
	for _, pa := range pathArray {
		if pa == "GOPATH" {
			gopath := os.Getenv("GOPATH")
			if gopath == "" {
				return "", fmt.Errorf("GOPATH environment is empty")
			} else {
				sanitizedArray = append(sanitizedArray, gopath)
			}
		} else {
			sanitizedArray = append(sanitizedArray, pa)
		}
	}
	return path.Join(sanitizedArray...), nil
}

type Parameters struct {
	model.In
	Env        env.Service
	Preference *config.Preference `optional:"true"`
}
