package impl

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func (s *Service) IsValid() (isValid bool) {
	return s.secret != ""
}

func (s *Service) Parse(token string) (claims map[string]interface{}, err error) {

	// *Code from https://pkg.go.dev/github.com/golang-jwt/jwt#example-Parse-Hmac

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return claims, nil
	} else {
		return nil, errors.New("token not valid")
	}

}

func (s *Service) Generate(userID int64, expiresAt time.Time, issuer string, claims map[string]interface{}) (token string, err error) {

	if userID <= 0 {
		return "", errors.New("user id in token generation cannot be <= 0")
	}

	if expiresAt.IsZero() {
		return "", errors.New("token expires cannot be zero")
	}

	if claims == nil {
		claims = map[string]interface{}{}
	}

	// Structured version of Claims Section, as referenced at
	// https://tools.ietf.org/html/rfc7519#section-4.1
	// See examples for how to use this with your own claim types
	// type StandardClaims struct {
	// 	Audience  string `json:"aud,omitempty"`
	// 	ExpiresAt int64  `json:"exp,omitempty"`
	// 	Id        string `json:"jti,omitempty"`
	// 	IssuedAt  int64  `json:"iat,omitempty"`
	// 	Issuer    string `json:"iss,omitempty"`
	// 	NotBefore int64  `json:"nbf,omitempty"`
	// 	Subject   string `json:"sub,omitempty"`
	// }

	claims["user_id"] = userID
	claims["exp"] = expiresAt.Unix()
	claims["iat"] = time.Now().Unix()
	claims["iss"] = issuer

	// *Code from https://pkg.go.dev/github.com/golang-jwt/jwt#example-New-Hmac

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	//
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"foo": "bar",
	// 	"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	// })
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	// Sign and get the complete encoded token as a string using the secret
	token, err = jwtToken.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return token, nil
}
