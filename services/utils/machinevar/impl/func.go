package impl

import (
	"fmt"
	"os"
)

func (s *Service) GetVar(key string) (result string, err error) {
	value := os.Getenv(key)
	if value == "" {
		return result, fmt.Errorf("environment variable with key [%v] is not found or empty", key)
	}
	return value, nil
}
