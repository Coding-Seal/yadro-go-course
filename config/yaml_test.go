package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	_, err := NewConfig("../tests/config/right.yaml")
	assert.NoError(t, err)
	_, err = NewConfig("../tests/config/no_file.yaml")
	assert.Error(t, err)
}
