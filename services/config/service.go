package config

import (
	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/services/env"
)

type Service struct {
	Server Server `yaml:"server"`
}

type Server struct {
	Port            int  `yaml:"port"`
	IsLoggingInJSON bool `yaml:"is_logging_in_json"`
}

var (
	GetDefaultGolangPath                          = func() []string { return []string{"GOPATH", "src"} }
	GetDefaultDevelopmentConfigFileRelativeToRepo = func() []string { return []string{"files", "var", "config", "ss.development.config.yaml"} }
	GetDefaultStagingConfigFileRelativeToRepo     = func() []string { return []string{"files", "var", "config", "ss.staging.config.yaml"} }
	GetDefaultProductionConfigFileRelativeToRepo  = func() []string { return []string{"files", "var", "config", "ss.production.config.yaml"} }

	GetDefaultFilePathUsingRepositoryPath = func(pathToRepo ...string) []string {
		result := functions.CombineStringArray(
			GetDefaultGolangPath(),
			append([]string{}, pathToRepo...),
			GetDefaultDevelopmentConfigFileRelativeToRepo())
		return result
	}
	GetServerDefaultFilePathUsingRepositoryPath = func(pathToRepo ...string) []string {
		result := functions.CombineStringArray(
			append([]string{}, pathToRepo...),
			GetDefaultDevelopmentConfigFileRelativeToRepo())
		return result
	}

	DEFAULT_FILE_PATH = map[string][]string{
		env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT: functions.CombineStringArray(GetDefaultGolangPath(), []string{"github.com", "gndw", "gank"}, GetDefaultDevelopmentConfigFileRelativeToRepo()),
		env.DEFAULT_ENV_NAME_ENV_STAGING:     append([]string{"/app"}, GetDefaultStagingConfigFileRelativeToRepo()...),
		env.DEFAULT_ENV_NAME_ENV_PRODUCTION:  append([]string{"/app"}, GetDefaultProductionConfigFileRelativeToRepo()...),
	}

	DEFAULT_CONFIG = Service{
		Server: Server{
			Port:            9000,
			IsLoggingInJSON: false,
		},
	}
)

type Preference struct {
	// currently application has default file path for default env. check DEFAULT_FILE_PATH
	// add your custom file path based on env here
	// example EnvFilePaths["my-custom-env"] = []string{ "/app","files","my-custom-config-file.yaml" }
	EnvFilePaths map[string][]string
}

func CreatePreference(preference Preference) func() (*Preference, error) {
	return func() (*Preference, error) { return &preference, nil }
}

type Content struct {
	Value []byte
}

func CreateDevelopmentPreference(configFileFolders ...string) func() (*Preference, error) {
	return CreatePreference(Preference{
		EnvFilePaths: map[string][]string{
			env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT: GetDefaultFilePathUsingRepositoryPath(configFileFolders...),
		},
	})
}
