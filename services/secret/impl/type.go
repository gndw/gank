package impl

import "github.com/gndw/gank/services/secret"

type Service struct {
	Database Database `json:"database"`
	Token    Token    `json:"token"`
}

func (m *Service) GetDatabase() secret.Database {
	if m != nil {
		return &m.Database
	}
	return nil
}

func (m *Service) GetToken() secret.Token {
	if m != nil {
		return &m.Token
	}
	return nil
}

type Database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"db_name"`
}

func (m *Database) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *Database) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *Database) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *Database) GetPort() int {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Database) GetDBName() string {
	if m != nil {
		return m.DBName
	}
	return ""
}

type Token struct {
	Key string `json:"key"`
}

func (m *Token) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}
