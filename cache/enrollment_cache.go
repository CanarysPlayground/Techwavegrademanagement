package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"techwave/models"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// EnrollmentCacheTTL is the time-to-live for cached enrollments (5 minutes)
	EnrollmentCacheTTL = 5 * time.Minute
	// EnrollmentCachePrefix is the prefix for enrollment cache keys
	EnrollmentCachePrefix = "enrollment:"
)

// EnrollmentCache provides Redis caching for enrollment data
type EnrollmentCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewEnrollmentCache creates a new enrollment cache instance
func NewEnrollmentCache(client *redis.Client) *EnrollmentCache {
	return &EnrollmentCache{
		client: client,
		ctx:    context.Background(),
	}
}

// Get retrieves an enrollment from cache
func (c *EnrollmentCache) Get(id string) (*models.Enrollment, error) {
	key := c.buildKey(id)
	
	data, err := c.client.Get(c.ctx, key).Bytes()
	if err == redis.Nil {
		// Cache miss
		return nil, nil
	}
	if err != nil {
		// Redis error - log but don't fail
		log.Printf("Redis Get error for key %s: %v", key, err)
		return nil, err
	}

	var enrollment models.Enrollment
	if err := json.Unmarshal(data, &enrollment); err != nil {
		log.Printf("Failed to unmarshal cached enrollment: %v", err)
		return nil, err
	}

	log.Printf("Cache HIT for enrollment ID: %s", id)
	return &enrollment, nil
}

// Set stores an enrollment in cache with TTL
func (c *EnrollmentCache) Set(enrollment *models.Enrollment) error {
	key := c.buildKey(enrollment.ID)
	
	data, err := json.Marshal(enrollment)
	if err != nil {
		log.Printf("Failed to marshal enrollment for caching: %v", err)
		return err
	}

	err = c.client.Set(c.ctx, key, data, EnrollmentCacheTTL).Err()
	if err != nil {
		log.Printf("Redis Set error for key %s: %v", key, err)
		return err
	}

	log.Printf("Cached enrollment ID: %s (TTL: %v)", enrollment.ID, EnrollmentCacheTTL)
	return nil
}

// Delete removes an enrollment from cache (for invalidation)
func (c *EnrollmentCache) Delete(id string) error {
	key := c.buildKey(id)
	
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		log.Printf("Redis Delete error for key %s: %v", key, err)
		return err
	}

	log.Printf("Cache invalidated for enrollment ID: %s", id)
	return nil
}

// buildKey constructs the Redis key for an enrollment
func (c *EnrollmentCache) buildKey(id string) string {
	return fmt.Sprintf("%s%s", EnrollmentCachePrefix, id)
}

// Ping checks if Redis connection is healthy
func (c *EnrollmentCache) Ping() error {
	return c.client.Ping(c.ctx).Err()
}

// GetStats returns basic cache statistics
func (c *EnrollmentCache) GetStats() (map[string]interface{}, error) {
	info, err := c.client.Info(c.ctx, "stats").Result()
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"info": info,
		"connected": c.client.Ping(c.ctx).Err() == nil,
	}, nil
}
