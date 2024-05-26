package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewComicService(t *testing.T) {
	assert.NotNil(t, NewComicService(newComicRepoMock()))
}
