package token

import "time"

type Service interface {
	IsValid() (isValid bool)
	Parse(token string) (claims map[string]interface{}, err error)
	Generate(userID int64, expiresAt time.Time, issuer string, claims map[string]interface{}) (token string, err error)
}
