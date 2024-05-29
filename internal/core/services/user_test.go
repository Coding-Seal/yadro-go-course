package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserService(t *testing.T) {
	assert.NotNil(t, NewUserService(newUserRepoMock()))
}
