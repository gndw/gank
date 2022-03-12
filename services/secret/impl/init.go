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
	"github.com/gndw/gank/services/secret"
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

func New(params Parameters) (data secret.Service, content secret.Content, err error) {

	data = secret.DEFAULT_SECRET

	internalData, err := PopulateDataFromPreference(params.Preference)
	if err != nil {
		return data, content, err
	}

	configPath, exist := internalData.eligibleEnvBasedFilePaths[params.Env.Get()]
	if exist {
		contentByte, err := PopulateDataFromSecretFilePath(params.Marshal, configPath.Path, &data)
		content.Value = contentByte
		if err != nil {
			if configPath.MustValid {
				return data, content, err
			} else {
				params.Log.Debugf("secret.service> failed to load default secret file for env[%v]. returning empty secret file", params.Env.Get())
			}
		} else {
			params.Log.Debugf("secret.service> successfully populate secret using file from %v", configPath.Path)
		}
	} else {
		params.Log.Debugln("secret.service> no secret file path. returning empty secret file")
	}

	return data, content, nil
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

func PopulateDataFromSecretFilePath(marshal marshal.Service, path string, target *secret.Service) (secretByte []byte, err error) {
	if path == "" {
		return secretByte, errors.New("secret file path cannot be empty")
	}
	secretByte, err = ioutil.ReadFile(path)
	if err != nil {
		return secretByte, fmt.Errorf("failed to read secret file in %v with err: %v", path, err)
	}
	err = marshal.JsonUnmarshal(secretByte, target)
	if err != nil {
		return secretByte, fmt.Errorf("failed to unmarshal secret file %v with err: %v", path, err)
	}
	return secretByte, nil
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
	Preference *secret.Preference `optional:"true"`
}
