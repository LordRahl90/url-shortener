package cache

import (
	"context"
	"errors"

	"shortener/domains/cache/storage"
)

var (
	_ storage.ICacheStore = (*MockCacheStore)(nil)

	errMockNotInitialized = errors.New("mock not initialized")
)

type MockCacheStore struct {
	FindFunc func(ctx context.Context, key string) (string, error)
	SaveFunc func(ctx context.Context, key string, value string) error
}

// Find implements storage.ICacheStore
func (m *MockCacheStore) Find(ctx context.Context, key string) (string, error) {
	if m.FindFunc == nil {
		return "", errMockNotInitialized
	}

	return m.FindFunc(ctx, key)
}

// Save implements storage.ICacheStore
func (m *MockCacheStore) Save(ctx context.Context, key string, value string) error {
	if m.SaveFunc == nil {
		return errMockNotInitialized
	}
	return m.SaveFunc(ctx, key, value)
}
