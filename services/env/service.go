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
	DEFAULT_ALLOWED_ENV_NAME         = []string{
		DEFAULT_ENV_NAME_ENV_DEVELOPMENT,
		DEFAULT_ENV_NAME_ENV_STAGING,
		DEFAULT_ENV_NAME_ENV_PRODUCTION,
	}

	DEFAULT_MACHINE_ENV_NAME = "APP_ENV"
)

type Preference struct {

	// default environment name if no environment is found
	// replacing DEFAULT_ENV_NAME_DEVELOPMENT
	DefaultEnv string

	// machine environment variable key to get environment
	// example : apps will find env name from machine environment variable with key : APP_ENV
	// you can change APP_ENV key to other name
	// replacing DEFAULT_MACHINE_ENV_NAME
	MachineEnvName string

	// currently application only allow development, staging, and production as env name
	// add your custom env name here
	AdditionalEnvs []string
}
