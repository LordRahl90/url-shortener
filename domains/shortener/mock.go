package shortener

import (
	"context"
	"errors"

	"shortener/domains/shortener/storage"
)

var (
	_ storage.IShortenerStorage = (*MockStorage)(nil)

	errMockNotInitialized = errors.New("mock not initialized")
)

// MockStorage a mock service for the storage
type MockStorage struct {
	CreateFunc      func(ctx context.Context, r *storage.Record) error
	FindByLongFunc  func(ctx context.Context, long string) (*storage.Record, error)
	FindByShortFunc func(ctx context.Context, short string) (*storage.Record, error)
}

// Create mocks storage create function
func (m *MockStorage) Create(ctx context.Context, r *storage.Record) error {
	if m.CreateFunc == nil {
		return errMockNotInitialized
	}
	return m.CreateFunc(ctx, r)
}

// FindByLong mocks storage FindByLong function
func (m *MockStorage) FindByLong(ctx context.Context, long string) (*storage.Record, error) {
	if m.FindByLongFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.FindByLongFunc(ctx, long)
}

// FindByShort mocks storage findByShort function
func (m *MockStorage) FindByShort(ctx context.Context, short string) (*storage.Record, error) {
	if m.FindByShortFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.FindByShortFunc(ctx, short)
}
