package impl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/secret"
	"github.com/gndw/gank/services/utils/log"
	"github.com/gndw/gank/services/utils/marshal"
)

type ConfigInternalData struct {
	eligibleEnvBasedFilePaths   map[string]ConfigPath
	eligibleEnvBasedMachineKeys map[string]ConfigPath
}

type ConfigPath struct {
	MustValid bool
	Value     string
}

func New(params Parameters) (data secret.Service, content secret.Content, err error) {

	data = secret.DEFAULT_SECRET

	internalData, err := PopulateDataFromPreference(params.Preference)
	if err != nil {
		return data, content, err
	}

	secretPath, exist := internalData.eligibleEnvBasedFilePaths[params.Env.Get()]
	if exist {
		contentByte, err := PopulateDataFromSecretFilePath(params.Marshal, secretPath.Value, &data)
		content.Value = contentByte
		if err != nil {
			if secretPath.MustValid {
				return data, content, err
			} else {
				params.Log.Debugf("secret.service> failed to load default secret file for env[%v]", params.Env.Get())
			}
		} else {
			params.Log.Debugf("secret.service> successfully populate secret using file from %v", secretPath.Value)
			return data, content, nil
		}
	} else {
		params.Log.Debugf("secret.service> no secret file path for this env [%v].", params.Env.Get())
	}

	machineKey, exist := internalData.eligibleEnvBasedMachineKeys[params.Env.Get()]
	if exist {
		contentByte, err := PopulateDataFromMachineEnvironmentVariable(params.Marshal, machineKey.Value, &data)
		content.Value = contentByte
		if err != nil {
			if machineKey.MustValid {
				return data, content, err
			} else {
				params.Log.Debugf("secret.service> failed to load machine env-var key [%v] for env[%v]", machineKey.Value, params.Env.Get())
			}
		} else {
			params.Log.Debugf("secret.service> successfully populate secret using machine env-var key %v", machineKey.Value)
			return data, content, nil
		}
	} else {
		params.Log.Debugf("secret.service> no secret machine env-var key for this env [%v].", params.Env.Get())
	}

	params.Log.Debugf("secret.service> no secret found for this env [%v]. Returning empty secret file", params.Env.Get())
	return data, content, nil
}

func PopulateDataFromPreference(pref *secret.Preference) (data ConfigInternalData, err error) {

	for envName, filePathFolders := range secret.DEFAULT_FILE_PATH {
		path, _ := GetPathFromArray(filePathFolders)
		if data.eligibleEnvBasedFilePaths == nil {
			data.eligibleEnvBasedFilePaths = make(map[string]ConfigPath)
		}
		data.eligibleEnvBasedFilePaths[envName] = ConfigPath{Value: path}
	}

	for envName, machineKey := range secret.DEFAULT_MACHINE_VAR {
		if data.eligibleEnvBasedMachineKeys == nil {
			data.eligibleEnvBasedMachineKeys = make(map[string]ConfigPath)
		}
		data.eligibleEnvBasedMachineKeys[envName] = ConfigPath{Value: machineKey}
	}

	if pref != nil {
		for envName, filePathFolders := range pref.EnvFilePaths {
			if !functions.IsAllNonEmpty(filePathFolders...) {
				return data, fmt.Errorf("secret file path for env [%v] cannot be empty : %v", envName, filePathFolders)
			}
			path, err := GetPathFromArray(filePathFolders)
			if err != nil {
				return data, err
			}
			data.eligibleEnvBasedFilePaths[envName] = ConfigPath{Value: path, MustValid: true}
		}

		for envName, machineVar := range pref.EnvMachineVar {
			if !functions.IsAllNonEmpty(machineVar) {
				return data, fmt.Errorf("secret machine env-var for env [%v] cannot be empty : %v", envName, machineVar)
			}
			data.eligibleEnvBasedMachineKeys[envName] = ConfigPath{Value: machineVar, MustValid: true}
		}
	}
	return data, nil
}

func PopulateDataFromSecretFilePath(marshal marshal.Service, path string, target *secret.Service) (content []byte, err error) {
	if path == "" {
		return content, errors.New("secret file path cannot be empty")
	}
	content, err = ioutil.ReadFile(path)
	if err != nil {
		return content, fmt.Errorf("failed to read secret file in %v with err: %v", path, err)
	}
	err = marshal.JsonUnmarshal(content, target)
	if err != nil {
		return content, fmt.Errorf("failed to unmarshal secret file %v with err: %v", path, err)
	}
	return content, nil
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

func PopulateDataFromMachineEnvironmentVariable(marshal marshal.Service, machineEnvKey string, target *secret.Service) (content []byte, err error) {
	if machineEnvKey == "" {
		return content, errors.New("machine env key cannot be empty")
	}
	str := os.Getenv(machineEnvKey)
	if str == "" {
		return content, fmt.Errorf("machine env key value [%v] not found", machineEnvKey)
	}
	err = marshal.JsonUnmarshal([]byte(str), target)
	if err != nil {
		return content, fmt.Errorf("failed to unmarshal secret value %v with err: %v", machineEnvKey, err)
	}
	return content, nil
}

type Parameters struct {
	model.In
	Env        env.Service
	Log        log.Service
	Marshal    marshal.Service
	Preference *secret.Preference `optional:"true"`
}
