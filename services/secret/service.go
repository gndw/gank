package secret

type Service interface {
	GetDatabase() Database
	GetToken() Token
}

type Database interface {
	GetUser() string
	GetPassword() string
	GetHost() string
	GetPort() int
	GetDBName() string
}

type Token interface {
	GetKey() string
}
