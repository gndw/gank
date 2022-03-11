package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/secret"
)

type ConfigInternalData struct {
	eligibleEnvBasedFilePaths map[string]ConfigPath // map[env]ConfigPath
}

type ConfigPath struct {
	MustValid bool
	Path      string
}

func New(params Parameters) (data secret.Service, err error) {

	internalData, err := PopulateDataFromPreference(params.Preference)
	if err != nil {
		return data, err
	}

	configPath, exist := internalData.eligibleEnvBasedFilePaths[params.Env.Get()]
	if exist {
		err := PopulateDataFromConfigFilePath(configPath.Path, data)
		if err != nil {
			if configPath.MustValid {
				return data, err
			}
		}
	}

	return data, nil
}

func PopulateDataFromPreference(pref *secret.Preference) (data ConfigInternalData, err error) {

	for envName, filePathFolders := range config.DEFAULT_FILE_PATH {
		path, _ := GetPathFromArray(filePathFolders)
		if data.eligibleEnvBasedFilePaths == nil {
			data.eligibleEnvBasedFilePaths = make(map[string]ConfigPath)
		}
		data.eligibleEnvBasedFilePaths[envName] = ConfigPath{Path: path}
	}

	if pref != nil {
		for envName, filePathFolders := range pref.EnvFilePaths {
			if !functions.IsAllNonEmpty(filePathFolders...) {
				return data, fmt.Errorf("config file path for env [%v] cannot be empty : %v", envName, filePathFolders)
			}
			path, err := GetPathFromArray(filePathFolders)
			if err != nil {
				return data, err
			}
			data.eligibleEnvBasedFilePaths[envName] = ConfigPath{Path: path, MustValid: true}
		}
	}
	return data, nil
}

func PopulateDataFromConfigFilePath(path string, target secret.Service) (err error) {
	if path == "" {
		return errors.New("config file path cannot be empty")
	}
	configByte, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file in %v with err: %v", path, err)
	}
	err = json.Unmarshal(configByte, &target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config file")
	}
	return nil
}

func GetPathFromArray(pathArray []string) (string, error) {

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
	Preference *secret.Preference `optional:"true"`
}
