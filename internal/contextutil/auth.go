package contextutil

import (
	"context"
)

type userIDKey struct{}

func UserID(ctx context.Context) int64 {
	id, ok := ctx.Value(userIDKey{}).(int64)
	if !ok {
		panic("Should have set userID")
	}

	return id
}

func WithUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

type isAdminKey struct{}

func IsAdmin(ctx context.Context) bool {
	isAdmin, ok := ctx.Value(isAdminKey{}).(bool)
	if !ok {
		panic("Should have set isAdmin")
	}

	return isAdmin
}

func WithIsAdmin(ctx context.Context, isAdmin bool) context.Context {
	return context.WithValue(ctx, isAdminKey{}, isAdmin)
}
