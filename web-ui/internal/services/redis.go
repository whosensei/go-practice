package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisService handles Redis operations
type RedisService struct {
	client *redis.Client
}

// NewRedisService creates a new Redis service instance
func NewRedisService() *RedisService {
	// Get Redis configuration from environment variables with defaults
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0 // Default DB

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password:     redisPassword,
		DB:           redisDB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	return &RedisService{
		client: client,
	}
}

// Connect establishes connection to Redis and verifies it
func (rs *RedisService) Connect(ctx context.Context) error {
	// Test connection
	_, err := rs.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return nil
}

// Close closes the Redis connection
func (rs *RedisService) Close() error {
	if rs.client != nil {
		return rs.client.Close()
	}
	return nil
}

// GetClient returns the underlying Redis client
func (rs *RedisService) GetClient() *redis.Client {
	return rs.client
}

// Set sets a key-value pair with optional expiration
func (rs *RedisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rs.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (rs *RedisService) Get(ctx context.Context, key string) (string, error) {
	return rs.client.Get(ctx, key).Result()
}

// Delete deletes one or more keys
func (rs *RedisService) Delete(ctx context.Context, keys ...string) error {
	return rs.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (rs *RedisService) Exists(ctx context.Context, keys ...string) (int64, error) {
	return rs.client.Exists(ctx, keys...).Result()
}

// Expire sets an expiration time on a key
func (rs *RedisService) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rs.client.Expire(ctx, key, expiration).Err()
}

// Increment increments a key's value
func (rs *RedisService) Increment(ctx context.Context, key string) (int64, error) {
	return rs.client.Incr(ctx, key).Result()
}

// Decrement decrements a key's value
func (rs *RedisService) Decrement(ctx context.Context, key string) (int64, error) {
	return rs.client.Decr(ctx, key).Result()
}

// IsConnected checks if Redis is connected
func (rs *RedisService) IsConnected(ctx context.Context) bool {
	_, err := rs.client.Ping(ctx).Result()
	return err == nil
}
