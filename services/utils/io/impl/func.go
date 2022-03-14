package impl

import (
	"io/ioutil"
)

func (s *Service) ReadFile(filepath string) (content []byte, err error) {
	return ioutil.ReadFile(filepath)
}
