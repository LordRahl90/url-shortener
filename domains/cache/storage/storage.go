package storage

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var _ ICacheStore = (*Redis)(nil)

// Redis implementation for cache storage
type Redis struct {
	Client *redis.Client
}

// New returns a new implementation for cache service
func New(client *redis.Client) ICacheStore {
	return &Redis{Client: client}
}

// Find implements finding a key from redis
func (*Redis) Find(ctx context.Context, key string) (any, error) {
	panic("unimplemented")
}

// Save implements saving a record to redis
func (*Redis) Save(ctx context.Context, key string, data any) error {
	panic("unimplemented")
}
