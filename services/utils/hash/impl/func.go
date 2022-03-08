package impl

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) IsPasswordMatch(ctx context.Context, password string, hash string) (err error) {

	if password == "" || hash == "" {
		return errors.New("hash or password cannot be empty")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password hash not matched")
	}
	return nil
}

func (s *Service) GetPasswordHash(ctx context.Context, password string) (hash string, err error) {

	if password == "" {
		return hash, errors.New("password to hash cannot be empty")
	}

	hashInByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return hash, fmt.Errorf("password hash not matched. err: %v", err)
	}

	return string(hashInByte), nil
}
