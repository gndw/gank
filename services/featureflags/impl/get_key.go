package impl

import (
	"fmt"
)

func (s *Service) GetKey(key string) (value string, err error) {
	if value, exist := s.FeatureFlags[key]; exist {
		return value, nil
	}
	return "", fmt.Errorf("feature flag with key [%v] is not found", key)
}
