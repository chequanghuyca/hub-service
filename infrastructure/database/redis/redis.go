package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient() (*RedisClient, error) {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	password := os.Getenv("REDIS_PASSWORD")

	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		if parsed, err := strconv.Atoi(dbStr); err == nil {
			db = parsed
		}
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	ctx := context.Background()

	// Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connection established successfully!")

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// Get retrieves a value from Redis by key
func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key doesn't exist
	}
	return val, err
}

// Set stores a value in Redis
func (r *RedisClient) Set(key string, value interface{}) error {
	return r.client.Set(r.ctx, key, value, 0).Err()
}

// SetWithExpiry stores a value with an expiration time
func (r *RedisClient) SetWithExpiry(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Delete removes a key from Redis
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (r *RedisClient) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// SetNX sets a key only if it doesn't exist (useful for deduplication)
func (r *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(r.ctx, key, value, expiration).Result()
}

// HealthCheck verifies Redis connection is alive
func (r *RedisClient) HealthCheck() error {
	_, err := r.client.Ping(r.ctx).Result()
	return err
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// GetClient returns the underlying Redis client for advanced operations
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// GetContext returns the context
func (r *RedisClient) GetContext() context.Context {
	return r.ctx
}
