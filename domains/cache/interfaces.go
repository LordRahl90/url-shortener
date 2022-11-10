package cache

import "context"

// ICacheService defines the expected construct for any cache service
type ICacheService interface {
	Save(ctx context.Context, key string, data any) error
	Find(ctx context.Context, key string) (any, error)
}
