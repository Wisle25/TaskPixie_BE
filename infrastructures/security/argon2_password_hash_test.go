package security_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/infrastructures/security"
)

func TestArgon2PasswordHash(t *testing.T) {
	argon2PasswordHash := security.NewArgon2()

	t.Run("Hash Password", func(t *testing.T) {
		// Arrange
		plainPassword := "plain_password"

		// Act
		encryptedPassword := argon2PasswordHash.Hash(plainPassword)

		// Assert
		assert.NotEqual(t, plainPassword, encryptedPassword, "encrypted password should not be the same as plain password")
		assert.IsType(t, "string", encryptedPassword, "encrypted password should be a string")
	})

	t.Run("Compare Password", func(t *testing.T) {
		plainPassword := "plain_password"

		// Act
		encryptedPassword := argon2PasswordHash.Hash(plainPassword)

		t.Run("Comparison shouldnt return error when its matched", func(t *testing.T) {
			// Action and Assert
			assert.NotPanics(t, func() {
				argon2PasswordHash.Compare(plainPassword, encryptedPassword)
			})
		})

		t.Run("Comparison should raise panic when its not matched", func(t *testing.T) {
			// Action
			assert.PanicsWithError(t, "Password is incorrect!", func() {
				argon2PasswordHash.Compare("anotherPassword", encryptedPassword)
			})
		})

	})
}
