package env

type Service interface {
	Get() (env string)
	IsDevelopment() (isDevelopment bool)
	IsStaging() (isStaging bool)
	IsProduction() (isProduction bool)
}
