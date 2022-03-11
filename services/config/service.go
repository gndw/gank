package config

import (
	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/services/env"
)

type Service struct {
	Token  Token
	Server Server
	Auth   Auth
}

type Token struct {
	Duration int64
}

type Server struct {
	Port int
}

type Auth struct {
	MinLengthNewAccountUsername int
	MaxLengthNewAccountUsername int
	MinLengthNewAccountPassword int
	MaxLengthNewAccountPassword int
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
