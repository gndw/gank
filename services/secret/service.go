package secret

import (
	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/services/env"
)

type Service struct {
	Database Database
	Token    Token
}

type Database struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

type Token struct {
	Key string
}

var (
	GetDefaultGolangPath                          = func() []string { return []string{"GOPATH", "src"} }
	GetDefaultDevelopmentSecretFileRelativeToRepo = func() []string { return []string{"files", "var", "secret", "ss.development.secret.yaml"} }
	GetDefaultStagingSecretFileRelativeToRepo     = func() []string { return []string{"files", "var", "secret", "ss.staging.secret.yaml"} }
	GetDefaultProductionSecretFileRelativeToRepo  = func() []string { return []string{"files", "var", "secret", "ss.production.secret.yaml"} }

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
