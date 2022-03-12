package secret

import (
	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/services/env"
)

type Service struct {
	Database Database `json:"database"`
	Token    Token    `json:"token"`
}

type Database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"db_name"`
}

type Token struct {
	Key string `json:"key"`
}

var (
	GetDefaultGolangPath                          = func() []string { return []string{"GOPATH", "src"} }
	GetDefaultDevelopmentSecretFileRelativeToRepo = func() []string { return []string{"files", "var", "secret", "ss.development.secret.json"} }
	GetDefaultStagingSecretFileRelativeToRepo     = func() []string { return []string{"files", "var", "secret", "ss.staging.secret.json"} }
	GetDefaultProductionSecretFileRelativeToRepo  = func() []string { return []string{"files", "var", "secret", "ss.production.secret.json"} }

	GetDefaultFilePathUsingRepositoryPath = func(pathToRepo ...string) []string {
		result := functions.CombineStringArray(
			GetDefaultGolangPath(),
			append([]string{}, pathToRepo...),
			GetDefaultDevelopmentSecretFileRelativeToRepo())
		return result
	}
	GetServerDefaultFilePathUsingRepositoryPath = func(pathToRepo ...string) []string {
		result := functions.CombineStringArray(
			append([]string{}, pathToRepo...),
			GetDefaultDevelopmentSecretFileRelativeToRepo())
		return result
	}

	DEFAULT_FILE_PATH = map[string][]string{
		env.DEFAULT_ENV_NAME_ENV_STAGING:    append([]string{"/app"}, GetDefaultStagingSecretFileRelativeToRepo()...),
		env.DEFAULT_ENV_NAME_ENV_PRODUCTION: append([]string{"/app"}, GetDefaultProductionSecretFileRelativeToRepo()...),
	}
)

type Preference struct {
	// currently application has default file path for default env. check DEFAULT_FILE_PATH
	// add your custom file path based on env here
	// example EnvFilePaths["my-custom-env"] = []string{ "/app","files","my-custom-secret-file.json" }
	EnvFilePaths map[string][]string
}

func CreatePreference(preference Preference) func() (*Preference, error) {
	return func() (*Preference, error) { return &preference, nil }
}
