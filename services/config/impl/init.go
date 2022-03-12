package impl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/utils/log"
	"github.com/gndw/gank/services/utils/marshal"
)

type ConfigInternalData struct {
	eligibleEnvBasedFilePaths map[string]ConfigPath // map[env]ConfigPath
}

type ConfigPath struct {
	MustValid bool
	Path      string
}

func New(params Parameters) (data config.Service, content config.Content, err error) {

	data = config.DEFAULT_CONFIG

	internalData, err := PopulateDataFromPreference(params.Preference)
	if err != nil {
		return data, content, err
	}

	configPath, exist := internalData.eligibleEnvBasedFilePaths[params.Env.Get()]
	if exist {
		contentByte, err := PopulateDataFromConfigFilePath(params.Marshal, configPath.Path, &data)
		content.Value = contentByte
		if err != nil {
			if configPath.MustValid {
				return data, content, err
			} else {
				params.Log.Debugf("config.service> failed to load default config file for env[%v]. returning empty config file", params.Env.Get())
			}
		} else {
			params.Log.Debugf("config.service> successfully populate config using file from %v", configPath.Path)
		}
	} else {
		params.Log.Debugln("config.service> no config file path. returning empty config file")
	}

	return data, content, nil
}

func PopulateDataFromPreference(pref *config.Preference) (data ConfigInternalData, err error) {

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

func PopulateDataFromConfigFilePath(marshal marshal.Service, path string, target *config.Service) (configByte []byte, err error) {
	if path == "" {
		return configByte, errors.New("config file path cannot be empty")
	}
	configByte, err = ioutil.ReadFile(path)
	if err != nil {
		return configByte, fmt.Errorf("failed to read config file in %v with err: %v", path, err)
	}
	err = marshal.YamlUnmarshal(configByte, target)
	if err != nil {
		return configByte, fmt.Errorf("failed to unmarshal config file %v with err: %v", path, err)
	}
	return configByte, nil
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
	Log        log.Service
	Marshal    marshal.Service
	Preference *config.Preference `optional:"true"`
}
