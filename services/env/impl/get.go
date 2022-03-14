package impl

import (
	"github.com/gndw/gank/services/env"
)

func (s *Service) Get() (env string) {
	return s.env
}

func (s *Service) IsDevelopment() (isDevelopment bool) {
	return s.env == env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT
}

func (s *Service) IsStaging() (isStaging bool) {
	return s.env == env.DEFAULT_ENV_NAME_ENV_STAGING
}

func (s *Service) IsProduction() (isProduction bool) {
	return s.env == env.DEFAULT_ENV_NAME_ENV_PRODUCTION
}

func (s *Service) IsReleaseLevel() (isReleaseLevel bool) {
	return s.isReleaseLevel
}
