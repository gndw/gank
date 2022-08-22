package config

import (
	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/services/env"
)

type Service struct {
	Server Server `yaml:"server"`
}

type Server struct {
	Port                          int    `yaml:"port"`
	IsLoggingInJSON               bool   `yaml:"is_logging_in_json"`
	IsLoggingFieldOnlyWhenVerbose bool   `yaml:"is_logging_field_only_when_verbose"`
	IsReturnDeveloperError        bool   `yaml:"is_return_developer_error"`
	DefaultMsgMaintenance         string `yaml:"default_msg_maintenance"`
	DefaultMsgUnauthorized        string `yaml:"default_msg_unauthorized"`
	DefaultMsgBadRequest          string `yaml:"default_msg_bad_request"`
	DefaultMsgInternalServerError string `yaml:"default_msg_internal_server_error"`

	AllowedOrigins   string `yaml:"allowed_origins"`
	AllowedMethods   string `yaml:"allowed_methods"`
	AllowedHeaders   string `yaml:"allowed_headers"`
	ExposedHeaders   string `yaml:"exposed_headers"`
	AllowCredentials bool   `yaml:"allow_credentials"`
	CacheMaxAge      int    `yaml:"cache_max_age"`
	SensitiveFields  string `yaml:"sensitive_fields"`
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
		env.DEFAULT_ENV_NAME_ENV_STAGING:     GetDefaultStagingConfigFileRelativeToRepo(),
		env.DEFAULT_ENV_NAME_ENV_PRODUCTION:  GetDefaultProductionConfigFileRelativeToRepo(),
	}

	DEFAULT_CONFIG = Service{
		Server: Server{
			Port:                          9000,
			IsLoggingInJSON:               false,
			IsLoggingFieldOnlyWhenVerbose: false,
			IsReturnDeveloperError:        true,
			DefaultMsgMaintenance:         "We are currently performing a scheduled maintenance. We will be back soon.",
			DefaultMsgUnauthorized:        "Sorry your session has timed-out. Please sign in again.",
			DefaultMsgBadRequest:          "We cannot proceed with your invalid request. Please modify and try again. (%v)",
			DefaultMsgInternalServerError: "We cannot proceed with your request at the moment. Please try again later. (%v)",
			AllowedOrigins:                "https://*,http://*",
			AllowedMethods:                "GET,POST,PUT,DELETE,OPTIONS",
			AllowedHeaders:                "Accept,Authorization,Content-Type,X-CSRF-Token,X-Request-ID",
			ExposedHeaders:                "Link",
			AllowCredentials:              false,
			CacheMaxAge:                   300,
			SensitiveFields:               "password",
		},
	}
)

type Preference struct {
	// currently application has default file path for default env. check DEFAULT_FILE_PATH
	// add your custom file path based on env here
	// example EnvFilePaths["my-custom-env"] = []string{ "files","var","my-custom-config-file.yaml" }
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
