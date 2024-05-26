package contextutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReqID(t *testing.T) {
	ctx := context.Background()

	assert.Panics(t, func() { ReqID(ctx) })

	ctx = WithReqID(ctx, 42)

	assert.NotPanics(t, func() { ReqID(ctx) })
	assert.Equal(t, 42, ReqID(ctx))
}
