package config

type Service interface {
	GetToken() Token
	GetServer() Server
	GetAuth() Auth
}

type Token interface {
	GetDuration() int64
}

type Server interface {
	GetPort() int
}

type Auth interface {
	GetMinLengthNewAccountUsername() int
	GetMaxLengthNewAccountUsername() int
	GetMinLengthNewAccountPassword() int
	GetMaxLengthNewAccountPassword() int
}
