package featureflags

type Service interface {
	GetKey(key string) (value string, err error)
}
