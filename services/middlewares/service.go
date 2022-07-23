package middlewares

import (
	"net/http"

	"github.com/gndw/gank/model"
)

type Service interface {
	GetInitializeMiddleware(f model.Middleware) http.HandlerFunc
	GetLoggerMiddleware(f model.Middleware) model.Middleware
	GetHttpMiddleware(f model.Middleware) model.Middleware
	GetRecovererMiddleware(f model.Middleware) model.Middleware
	GetRequestIDMiddleware(f model.Middleware) model.Middleware
	GetDefault() []func(m model.Middleware) model.Middleware
	GetDefaultWith(middlewares ...func(m model.Middleware) model.Middleware) []func(m model.Middleware) model.Middleware
}

type Auth interface {
	GetAuthMiddleware(f model.Middleware) model.Middleware
	GetAuthorizationHeader(r *http.Request) (isExist bool, headerStr string)
	IsBearerAuthentication(headerStr string) (isBearer bool, token string)
	ParseToken(token string, keys []string, config AuthParseTokenConfig) (mapOfClaimByKeys map[string]interface{}, err error)
}

type AuthParseTokenConfig struct {
	CheckClaimPolicy AuthCheckClaimPolicy
}

type AuthCheckClaimPolicy string

var IGNORE_MISSING_KEYS AuthCheckClaimPolicy = "ignore_missing_keys"
var MUST_VALID_KEYS AuthCheckClaimPolicy = "must_valid_keys"

var DEFAULT_CLAIM_KEY_USER_ID = "user_id"
