package generator_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/infrastructures/generator"
	"testing"
)

func TestUUIDGenerator(t *testing.T) {
	// Arrange
	var id string
	idGenerator := generator.NewUUIDGenerator()

	// Actions and Assert
	assert.NotPanics(t, func() {
		id = idGenerator.Generate()
	})
	assert.NotNil(t, id)
}
