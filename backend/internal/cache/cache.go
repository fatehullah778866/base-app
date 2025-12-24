package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

type inMemoryCache struct {
	store map[string]*cacheEntry
	logger *zap.Logger
}

type cacheEntry struct {
	value     []byte
	expiresAt time.Time
}

func NewInMemoryCache(logger *zap.Logger) Cache {
	c := &inMemoryCache{
		store:  make(map[string]*cacheEntry),
		logger: logger,
	}
	
	// Cleanup expired entries
	go c.cleanup()
	
	return c
}

func (c *inMemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		now := time.Now()
		for key, entry := range c.store {
			if now.After(entry.expiresAt) {
				delete(c.store, key)
			}
		}
	}
}

func (c *inMemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	entry, exists := c.store[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	
	if time.Now().After(entry.expiresAt) {
		delete(c.store, key)
		return nil, fmt.Errorf("key expired: %s", key)
	}
	
	return entry.value, nil
}

func (c *inMemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	c.store[key] = &cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

func (c *inMemoryCache) Delete(ctx context.Context, key string) error {
	delete(c.store, key)
	return nil
}

func (c *inMemoryCache) Clear(ctx context.Context) error {
	c.store = make(map[string]*cacheEntry)
	return nil
}

// Helper functions for common types
func GetJSON[T any](cache Cache, ctx context.Context, key string) (*T, error) {
	data, err := cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}
	
	return &value, nil
}

func SetJSON[T any](cache Cache, ctx context.Context, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return cache.Set(ctx, key, data, ttl)
}

