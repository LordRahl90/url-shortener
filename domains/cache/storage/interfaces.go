package storage

import "context"

// ICacheStore describes the expected features of cache service
type ICacheStore interface {
	Save(ctx context.Context, key, value string) error
	Find(ctx context.Context, key string) (string, error)
}
