package config

import (
	"github.com/gndw/gank/functions"
)

type Service interface {
	GetToken() Token
	GetServer() Server
	GetAuth() Auth
}

type Token interface {
	GetDuration() int64
}

type Server interface {
	GetPort() int
}

type Auth interface {
	GetMinLengthNewAccountUsername() int
	GetMaxLengthNewAccountUsername() int
	GetMinLengthNewAccountPassword() int
	GetMaxLengthNewAccountPassword() int
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

	DEFAULT_FILE_PATH_ON_DEVELOPMENT_SERVER = functions.CombineStringArray(GetDefaultGolangPath(), []string{"github.com", "gndw", "gank"}, GetDefaultDevelopmentConfigFileRelativeToRepo())
	DEFAULT_FILE_PATH_ON_STAGING_SERVER     = append([]string{"/app"}, GetDefaultStagingConfigFileRelativeToRepo()...)
	DEFAULT_FILE_PATH_ON_PRODUCTION_SERVER  = append([]string{"/app"}, GetDefaultProductionConfigFileRelativeToRepo()...)
)

type Preference struct {
	FilePathDevelopment    []string
	FilePathStaging        []string
	FilePathProduction     []string
	AdditionalEnvFilePaths map[string][]string // file path based on custom environment
}

func CreatePreference(preference Preference) func() (*Preference, error) {
	return func() (*Preference, error) { return &preference, nil }
}
