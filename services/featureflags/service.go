package featureflags

type Service interface {
	GetKey(key string) (value string, err error)
	GetKeyWithDefault(key string, def string) (value string)
	GetBoolean(key string) (value bool, err error)
	GetBooleanWithDefault(key string, def bool) (value bool)
	GetInt64(key string) (value int64, err error)
	GetInt64WithDefault(key string, def int64) (value int64)
	GetArrayOfInt64(key string) (value []int64, err error)
	GetArrayOfInt64WithDefault(key string, def []int64) (value []int64)
}
