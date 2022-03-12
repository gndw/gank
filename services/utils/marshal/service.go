package marshal

type Service interface {
	JsonMarshal(v interface{}) (result []byte, err error)
	JsonUnmarshal(data []byte, v interface{}) (err error)
	YamlMarshal(v interface{}) (result []byte, err error)
	YamlUnmarshal(data []byte, v interface{}) (err error)
}
