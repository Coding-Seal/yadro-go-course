package contextutil

import (
	"context"
	"errors"
)

type requestIDKey struct{}

func ReqID(ctx context.Context) (int, error) {
	id, ok := ctx.Value(requestIDKey{}).(int)
	if !ok {
		return 0, errors.New("invalid request id")
	}

	return id, nil
}

func WithReqID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}
