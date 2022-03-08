package impl

import (
	"github.com/gndw/gank/services/config"
)

type Service struct {
	Token  Token  `yaml:"token"`
	Server Server `yaml:"server"`
	Auth   Auth   `yaml:"auth"`
}

func (m *Service) GetToken() config.Token {
	if m != nil {
		return &m.Token
	}
	return nil
}
func (m *Service) GetServer() config.Server {
	if m != nil {
		return &m.Server
	}
	return nil
}
func (m *Service) GetAuth() config.Auth {
	if m != nil {
		return &m.Auth
	}
	return nil
}

type Token struct {
	Duration int64 `yaml:"duration"`
}

func (m *Token) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

type Server struct {
	Port int `yaml:"port"`
}

func (m *Server) GetPort() int {
	if m != nil {
		return m.Port
	}
	return 0
}

type Auth struct {
	MinLengthNewAccountUsername int `yaml:"min_length_new_account_username"`
	MaxLengthNewAccountUsername int `yaml:"max_length_new_account_username"`
	MinLengthNewAccountPassword int `yaml:"min_length_new_account_password"`
	MaxLengthNewAccountPassword int `yaml:"max_length_new_account_password"`
}

func (m *Auth) GetMinLengthNewAccountUsername() int {
	if m != nil {
		return m.MinLengthNewAccountUsername
	}
	return 0
}

func (m *Auth) GetMaxLengthNewAccountUsername() int {
	if m != nil {
		return m.MaxLengthNewAccountUsername
	}
	return 0
}

func (m *Auth) GetMinLengthNewAccountPassword() int {
	if m != nil {
		return m.MinLengthNewAccountPassword
	}
	return 0
}

func (m *Auth) GetMaxLengthNewAccountPassword() int {
	if m != nil {
		return m.MaxLengthNewAccountPassword
	}
	return 0
}
