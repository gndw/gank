package env

type Service interface {
	Get() (env string)
	IsDevelopment() (isDevelopment bool)
	IsStaging() (isStaging bool)
	IsProduction() (isProduction bool)
}

var (
	DEFAULT_ENV_NAME_ENV_DEVELOPMENT = "development"
	DEFAULT_ENV_NAME_ENV_STAGING     = "staging"
	DEFAULT_ENV_NAME_ENV_PRODUCTION  = "production"

	DEFAULT_FLAG_NAME_ENV = "env"

	DEFAULT_MACHINE_ENV_NAME = "APP_ENV"
)

type Preference struct {
	DefaultEnv     string // replacing DEFAULT_ENV_NAME_DEVELOPMENT
	FlagNameEnv    string // replacing DEFAULT_FLAG_NAME_ENV
	MachineEnvName string // replacing DEFAULT_MACHINE_ENV_NAME
	AdditionalEnvs []string
}
