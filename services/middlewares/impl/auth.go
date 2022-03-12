package impl

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gndw/gank/constant"
	"github.com/gndw/gank/errorsg"
	"github.com/gndw/gank/model"
)

func (s *Service) GetAuthMiddleware(isActivateAuth bool, f model.Middleware) model.Middleware {
	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {

		if isActivateAuth {
			ctx, err = s.ValidateAuthFromHeader(ctx, r)
			if err != nil {
				return nil, errorsg.WithOptions(err, errorsg.WithStatusCode(http.StatusUnauthorized))
			}
		}
		return f(ctx, rw, r)
	}
}

func (s *Service) ValidateAuthFromHeader(ctx context.Context, r *http.Request) (resultCtx context.Context, err error) {

	authHeader := r.Header.Get("Authorization")

	// Expecting auth header with format 'Bearer <token>'
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {

		// prepare token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := s.tokenService.Parse(tokenString)
		if err != nil {
			return ctx, err
		}

		// get jwt claims
		userIDFromClaim, exist := claims["user_id"]
		if !exist {
			return ctx, errors.New("token has no user_id claims")
		}

		// number in jwt token claims, always have float64 type data
		userID_Float, ok := userIDFromClaim.(float64)
		if !ok {
			return ctx, errors.New("invalid user_id claims")
		}

		// make sure user ID injected to context is int64
		resultCtx = context.WithValue(ctx, constant.ContextKeyUserID, int64(userID_Float))

		// success
		return resultCtx, nil

	} else {
		return ctx, errors.New("no valid authorization found")
	}
}

func (s *Service) IsAuthMiddlewareValid() (isValid bool) {
	return s.tokenService.IsValid()
}
