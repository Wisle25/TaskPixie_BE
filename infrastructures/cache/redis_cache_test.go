package cache_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/cache"
	"github.com/wisle25/be-template/infrastructures/services"
	"testing"
	"time"
)

func TestRedisCache(t *testing.T) {
	// Load configuration
	config := commons.LoadConfig("../../")

	// Connect to Redis
	redis := services.ConnectRedis(config)

	// Create RedisCache instance
	redisCache := cache.NewRedisCache(redis)
	ctx := context.TODO()

	t.Run("SetCache", func(t *testing.T) {
		// Arrange
		key := "test-key"
		value := "test-value"
		expiration := time.Minute

		// Act
		redisCache.SetCache(key, value, expiration)

		// Assert
		result, err := redis.Get(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("GetCache", func(t *testing.T) {
		t.Run("Should return nil when value is not found", func(t *testing.T) {
			// Act
			result := redisCache.GetCache("non-exists-key")

			// Assert
			assert.Nil(t, result)
		})

		t.Run("should get a value from Redis", func(t *testing.T) {
			// Arrange
			key := "test-key"
			expectedValue := "test-value"
			redis.Set(ctx, key, expectedValue, time.Minute)

			// Act
			result := redisCache.GetCache(key)

			// Assert
			assert.Equal(t, expectedValue, result)
		})
	})

	t.Run("DelCache", func(t *testing.T) {
		// Arrange
		key := "test-key"
		value := "test-value"

		redis.Set(ctx, key, value, time.Minute)

		// Act
		redisCache.DeleteCache(key)

		// Assert
		assert.Nil(t, redisCache.GetCache(key))
	})
}
