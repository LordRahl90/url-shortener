package cache

import (
	"context"

	"shortener/domains/cache/storage"
)

var _ ICacheService = (*CacheService)(nil)

// CacheService implements a concrete cache service class
type CacheService struct {
	store storage.ICacheStore
}

// New returns a new instance of cache store
func New(store storage.ICacheStore) ICacheService {
	return &CacheService{
		store: store,
	}
}

// Find implements finding an item in the cache
func (cs *CacheService) Find(ctx context.Context, key string) (string, error) {
	return cs.store.Find(ctx, key)
}

// Save implements keeping an item in the cache
func (cs *CacheService) Save(ctx context.Context, key string, data string) error {
	return cs.store.Save(ctx, key, data)
}
