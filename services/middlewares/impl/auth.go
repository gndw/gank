package impl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gndw/gank/contextg"
	"github.com/gndw/gank/errorsg"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/middlewares"
	"github.com/gndw/gank/services/utils/token"
)

type Auth struct {
	tokenService  token.Service
	configService config.Service
}

var DEFAULT_CLAIM_KEY_USER_ID = "user_id"

func NewAuth(token token.Service, config config.Service) (middlewares.Auth, error) {
	auth := &Auth{
		tokenService:  token,
		configService: config,
	}
	return auth, nil
}

func (s *Auth) GetAuthMiddleware(f model.Middleware) model.Middleware {
	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {

		ctx, err = s.validateAuthFromHeader(ctx, r)
		if err != nil {
			return nil, errorsg.WithOptions(err,
				errorsg.WithHttpStatusCode(http.StatusUnauthorized),
				errorsg.WithType(errorsg.ErrorTypeBadRequest),
				errorsg.WithPrettyMessage(s.configService.Server.DefaultMsgUnauthorized))
		}

		return f(ctx, rw, r)
	}
}

func (s *Auth) validateAuthFromHeader(ctx context.Context, r *http.Request) (resultCtx context.Context, err error) {

	isExist, headerStr := s.GetAuthorizationHeader(r)
	isBearer, token := s.IsBearerAuthentication(headerStr)

	if isExist && isBearer {

		mapOfClaimByKeys, err := s.ParseToken(token, []string{DEFAULT_CLAIM_KEY_USER_ID}, AuthParseTokenConfig{CheckClaimPolicy: MUST_VALID_KEYS})
		if err != nil {
			return ctx, err
		}

		userIDClaim, exist := mapOfClaimByKeys[DEFAULT_CLAIM_KEY_USER_ID]
		if exist {

			userID, err := s.ConvertClaimToInt64(userIDClaim)
			if err != nil {
				return ctx, err
			}

			resultCtx = contextg.WithUserID(ctx, userID)

		} else {
			resultCtx = ctx
		}

		// success
		return resultCtx, nil

	} else {
		return ctx, errors.New("no valid authorization found")
	}
}

func (s *Auth) GetAuthorizationHeader(r *http.Request) (isExist bool, headerStr string) {
	authHeader := r.Header.Get("Authorization")
	return authHeader != "", authHeader
}

func (s *Auth) IsBearerAuthentication(headerStr string) (isBearer bool, token string) {
	return strings.HasPrefix(headerStr, "Bearer "), strings.TrimPrefix(headerStr, "Bearer ")
}

type AuthParseTokenConfig struct {
	CheckClaimPolicy AuthCheckClaimPolicy
}

type AuthCheckClaimPolicy string

var IGNORE_MISSING_KEYS AuthCheckClaimPolicy = "ignore_missing_keys"
var MUST_VALID_KEYS AuthCheckClaimPolicy = "must_valid_keys"

func (s *Auth) ParseToken(token string, keys []string, config AuthParseTokenConfig) (mapOfClaimByKeys map[string]interface{}, err error) {

	// parse token
	claims, err := s.tokenService.Parse(token)
	if err != nil {
		return mapOfClaimByKeys, err
	}

	// default config
	if config.CheckClaimPolicy == "" {
		config.CheckClaimPolicy = IGNORE_MISSING_KEYS
	}

	// create map for result
	mapOfClaimByKeys = make(map[string]interface{})

	// get jwt claims
	for _, key := range keys {
		value, exist := claims[key]
		if exist {
			mapOfClaimByKeys[key] = value
		} else {
			switch config.CheckClaimPolicy {
			case MUST_VALID_KEYS:
				return mapOfClaimByKeys, fmt.Errorf("token has no %v claims", key)
			case IGNORE_MISSING_KEYS:
				// do nothing
			default:
				// do nothing
			}
		}
	}

	return mapOfClaimByKeys, nil
}

func (s *Auth) ConvertClaimToInt64(claim interface{}) (value int64, err error) {

	// number in jwt token claims, always have float64 type data
	claimInFloat64, ok := claim.(float64)
	if !ok {
		return value, errors.New("invalid user_id claims")
	}

	return int64(claimInFloat64), nil
}
