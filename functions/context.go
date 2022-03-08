package functions

import (
	"context"
	"errors"

	"github.com/gndw/gank/constant"
)

func GetUserIDFromContext(ctx context.Context) (userID int64, err error) {
	userID, ok := ctx.Value(constant.ContextKeyUserID).(int64)
	if !ok || userID == 0 {
		return userID, errors.New("user id not found")
	}
	return userID, nil
}
