package contextutil

import (
	"context"
)

type requestIDKey struct{}

func ReqID(ctx context.Context) int {
	id, ok := ctx.Value(requestIDKey{}).(int)
	if !ok {
		panic("Should have set request id")
	}

	return id
}

func WithReqID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}
