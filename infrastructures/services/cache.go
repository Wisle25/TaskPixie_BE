package services

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/wisle25/be-template/commons"
	"log"
)

// ConnectRedis initializes a connection to the Redis server using the provided configuration.
func ConnectRedis(config *commons.Config) *redis.Client {
	ctx := context.TODO()

	RedisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
	})

	// Ping the Redis server to ensure the connection is established.
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	log.Println("Successfully connected to Redis Client")

	return RedisClient
}
