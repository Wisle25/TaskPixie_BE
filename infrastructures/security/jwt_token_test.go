package security_test

import (
	"github.com/wisle25/be-template/infrastructures/generator"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/security"
)

func TestJWTToken(t *testing.T) {
	config := commons.LoadConfig("../..")
	uuidGenerator := generator.NewUUIDGenerator()

	jt := security.NewJwtToken(uuidGenerator)

	t.Run("Create JWT Token", func(t *testing.T) {
		t.Run("ValidTokenCreation", func(t *testing.T) {
			// Arrange
			userID := "test-user-id"
			ttl := time.Hour * 1

			// Action
			tokenDetail := jt.CreateToken(userID, ttl, config.AccessTokenPrivateKey)

			// Assert
			assert.NotNil(t, tokenDetail)
			assert.NotEmpty(t, tokenDetail.Token)
			assert.Equal(t, userID, tokenDetail.UserId)
			assert.WithinDuration(t, time.Now().Add(ttl), time.Unix(tokenDetail.ExpiresIn, 0), time.Minute)
			assert.NotEmpty(t, tokenDetail.TokenId)
		})

		t.Run("InvalidPrivateKey", func(t *testing.T) {
			// Arrange
			userID := "test-user-id"
			ttl := time.Hour * 1
			invalidPrivateKey := "invalid-private-key"

			// Action

			// Assert
			assert.Panics(t, func() {
				tokenDetail := jt.CreateToken(userID, ttl, invalidPrivateKey)
				assert.Nil(t, tokenDetail)
			})
		})

	})

	t.Run("Valiate JWT Token", func(t *testing.T) {
		t.Run("ValidToken", func(t *testing.T) {
			// Arrange
			userID := "test-user-id"
			ttl := time.Hour * 1

			// Action
			tokenDetail := jt.CreateToken(userID, ttl, config.AccessTokenPrivateKey)

			// Action
			validatedTokenDetail := jt.ValidateToken(tokenDetail.Token, config.AccessTokenPublicKey)

			// Assert
			require.NotNil(t, validatedTokenDetail)
			assert.Equal(t, tokenDetail.TokenId, validatedTokenDetail.TokenId)
			assert.Equal(t, tokenDetail.UserId, validatedTokenDetail.UserId)
		})

		t.Run("InvalidToken", func(t *testing.T) {
			// Arrange
			invalidToken := "invalid-token"

			// Action and Assert
			assert.Panics(t, func() {
				validatedTokenDetail := jt.ValidateToken(invalidToken, config.AccessTokenPublicKey)
				assert.Nil(t, validatedTokenDetail)
			})
		})

		t.Run("EmptyToken", func(t *testing.T) {
			// Arrange
			emptyToken := ""

			// Action and Assert
			assert.Panics(t, func() {
				validatedTokenDetail := jt.ValidateToken(emptyToken, config.AccessTokenPublicKey)
				assert.Nil(t, validatedTokenDetail)
			})
		})

		t.Run("InvalidPublicKey", func(t *testing.T) {
			// Arrange
			userID := "test-user-id"
			ttl := time.Hour * 1
			invalidPublicKey := "invalid-public-key"

			tokenDetail := jt.CreateToken(userID, ttl, config.AccessTokenPrivateKey)

			// Action and Assert
			assert.Panics(t, func() {
				validatedTokenDetail := jt.ValidateToken(tokenDetail.Token, invalidPublicKey)
				assert.Nil(t, validatedTokenDetail)
			})
		})
	})
}
