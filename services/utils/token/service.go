package token

import "time"

type Service interface {
	Parse(token string) (claims map[string]interface{}, err error)
	Generate(userID int64, expiresAt time.Time, issuer string, claims map[string]interface{}) (token string, err error)
}
