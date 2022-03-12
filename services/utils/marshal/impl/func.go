package impl

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

func (s *Service) JsonMarshal(v interface{}) (result []byte, err error) {
	return json.Marshal(v)
}

func (s *Service) JsonUnmarshal(data []byte, v interface{}) (err error) {
	return json.Unmarshal(data, v)
}

func (s *Service) YamlMarshal(v interface{}) (result []byte, err error) {
	return yaml.Marshal(v)
}

func (s *Service) YamlUnmarshal(data []byte, v interface{}) (err error) {
	return yaml.Unmarshal(data, v)
}
