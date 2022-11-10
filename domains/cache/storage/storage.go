package storage

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var _ ICacheStore = (*CaceheStore)(nil)

// Redis implementation for cache storage
type CaceheStore struct {
	client *redis.Client
}

// New returns a new implementation for cache service
func New(client *redis.Client) ICacheStore {
	return &CaceheStore{client: client}
}

// Find implements finding a key from redis
func (cs *CaceheStore) Find(ctx context.Context, key string) (string, error) {
	return cs.client.Get(ctx, key).Result()
}

// Save implements saving a record to redis
func (cs *CaceheStore) Save(ctx context.Context, key, value string) error {
	return cs.client.Set(ctx, key, value, 0).Err() // never expire
}
