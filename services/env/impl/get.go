package impl

import (
	"github.com/gndw/gank/constant"
)

func (s *Service) Get() (env string) {
	return s.env
}

func (s *Service) IsDevelopment() (isDevelopment bool) {
	return s.env == constant.ENV_DEVELOPMENT
}

func (s *Service) IsStaging() (isStaging bool) {
	return s.env == constant.ENV_STAGING
}

func (s *Service) IsProduction() (isProduction bool) {
	return s.env == constant.ENV_PRODUCTION
}
