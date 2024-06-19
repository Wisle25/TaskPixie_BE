package commons_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/commons"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Arrange
	var cfg *commons.Config

	// Actions and Assert
	assert.NotPanics(t, func() {
		cfg = commons.LoadConfig("..")
	})
	assert.NotNil(t, cfg)
	assert.Equal(t, "8000", cfg.ServerPort)
	assert.Equal(t, "test", cfg.AppEnv)
}
