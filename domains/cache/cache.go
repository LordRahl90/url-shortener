package cache

import (
	"context"
	"shortener/domains/cache/storage"
)

var _ ICacheService = (*CacheService)(nil)

// CacheService implements a concrete cache service class
type CacheService struct {
	Storage storage.ICacheStore
}

// New returns a new instance of cache store
func New(store storage.ICacheStore) ICacheService {
	return &CacheService{
		Storage: store,
	}
}

// Find implements finding an item in the cache
func (*CacheService) Find(ctx context.Context, key string) (any, error) {
	panic("unimplemented")
}

// Save implements keeping an item in the cache
func (*CacheService) Save(ctx context.Context, key string, data any) error {
	panic("unimplemented")
}
