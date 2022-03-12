package curruser

import (
	"context"
	"errors"
)

func CurrentUser(ctx context.Context) (uint64, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if userID, ok := ctx.Value(CurrentUserKey).(uint64); ok {
		return userID, nil
	}

	return 0, errors.New("unauthorized")
}

func SetCurrentUser(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, CurrentUserKey, userID)
}
