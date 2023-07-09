package impl

import (
	"fmt"
	"strconv"
	"strings"
)

func (s *Service) GetKey(key string) (value string, err error) {
	if value, exist := s.FeatureFlags[key]; exist {
		return value, nil
	}
	return "", fmt.Errorf("feature flag with key [%v] is not found", key)
}

func (s *Service) GetKeyWithDefault(key string, def string) (value string) {
	v, err := s.GetKey(key)
	if err == nil {
		return v
	}
	return def
}

func (s *Service) GetBoolean(key string) (value bool, err error) {
	v, err := s.GetKey(key)
	if err == nil {
		b, parseErr := strconv.ParseBool(v)
		if parseErr == nil {
			return b, nil
		}
		return false, fmt.Errorf("feature flag with key [%v] is not parsable to boolean", key)
	}
	return false, err
}

func (s *Service) GetBooleanWithDefault(key string, def bool) (value bool) {
	v, err := s.GetBoolean(key)
	if err == nil {
		return v
	}
	return def
}

func (s *Service) GetInt64(key string) (value int64, err error) {
	v, err := s.GetKey(key)
	if err == nil {
		b, parseErr := strconv.ParseInt(v, 10, 64)
		if parseErr == nil {
			return b, nil
		}
		return 0, fmt.Errorf("feature flag with key [%v] is not parsable to int64", key)
	}
	return 0, err
}

func (s *Service) GetInt64WithDefault(key string, def int64) (value int64) {
	v, err := s.GetInt64(key)
	if err == nil {
		return v
	}
	return def
}

func (s *Service) GetArrayOfInt64(key string) (value []int64, err error) {
	v, err := s.GetKey(key)
	if err == nil {
		valueStr := strings.Split(v, ",")
		for _, vStr := range valueStr {
			b, parseErr := strconv.ParseInt(vStr, 10, 64)
			if parseErr == nil {
				value = append(value, b)
			} else {
				fmt.Printf("[GetArrayOfInt64][feature-flags] error parsing %v from key %v", vStr, key)
			}
		}
		return value, nil
	}
	return []int64{}, err
}

func (s *Service) GetArrayOfInt64WithDefault(key string, def []int64) (value []int64) {
	v, err := s.GetArrayOfInt64(key)
	if err == nil {
		return v
	}
	return def
}
