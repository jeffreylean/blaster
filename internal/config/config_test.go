package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	err := os.Setenv("BLASTER_CONF", "file:///sample_config.yaml")
	assert.NoError(t, err)
	assert.NotPanics(t, func() {
		cfg := LoadOrPanic()
		assert.NotNil(t, cfg)
	})
}
