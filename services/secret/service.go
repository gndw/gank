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
		env.DEFAULT_ENV_NAME_ENV_STAGING:    GetDefaultStagingSecretFileRelativeToRepo(),
		env.DEFAULT_ENV_NAME_ENV_PRODUCTION: GetDefaultProductionSecretFileRelativeToRepo(),
	}

	DEFAULT_MACHINE_VAR = map[string]string{
		env.DEFAULT_ENV_NAME_ENV_STAGING:    DEFAULT_MACHINE_SECRET_NAME,
		env.DEFAULT_ENV_NAME_ENV_PRODUCTION: DEFAULT_MACHINE_SECRET_NAME,
	}

	DEFAULT_SECRET = Service{}

	DEFAULT_MACHINE_SECRET_NAME = "APP_SECRET"
)

type Preference struct {
	// currently application has default file path for default env. check DEFAULT_FILE_PATH
	// add your custom file path based on env here
	// example EnvFilePaths["my-custom-env"] = []string{ "files","var","my-custom-secret-file.json" }
	EnvFilePaths map[string][]string

	// currently application will check machine environment variable if file is not found as a file.
	// add your custom machine var based on env here
	// example EnvMachineVar["my-custom-env"] = APP_CUSTOM_SECRET
	EnvMachineVar map[string]string
}

func CreatePreference(preference Preference) func() (*Preference, error) {
	return func() (*Preference, error) { return &preference, nil }
}

type Content struct {
	Value []byte
}

func CreateDevelopmentPreference(secretFileFolders ...string) func() (*Preference, error) {
	return CreatePreference(Preference{
		EnvFilePaths: map[string][]string{
			env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT: GetDefaultFilePathUsingRepositoryPath(secretFileFolders...),
		},
	})
}
