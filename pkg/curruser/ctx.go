package curruser

import "context"

func CurrentUser(ctx context.Context) (uint64, context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	if userID, ok := ctx.Value(CurrentUserKey).(uint64); ok {
		return userID, ctx
	}

	return 0, ctx
}

func SetCurrentUser(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, CurrentUserKey, userID)
}
