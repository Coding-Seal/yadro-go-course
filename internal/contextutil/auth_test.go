package contextutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAdmin(t *testing.T) {
	ctx := context.Background()

	assert.Panics(t, func() { IsAdmin(ctx) })

	ctx = WithIsAdmin(ctx, true)

	assert.NotPanics(t, func() { IsAdmin(ctx) })
	assert.True(t, IsAdmin(ctx))
}

func TestUserID(t *testing.T) {
	ctx := context.Background()

	assert.Panics(t, func() { UserID(ctx) })

	ctx = WithUserID(ctx, 42)

	assert.NotPanics(t, func() { UserID(ctx) })
	assert.Equal(t, int64(42), UserID(ctx))
}
