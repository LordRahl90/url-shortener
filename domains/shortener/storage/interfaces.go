package storage

import "context"

// IShortenerStorage contract for shortener storage
type IShortenerStorage interface {
	Create(ctx context.Context, r *Record) error
	FindByShort(ctx context.Context, short string) (*Record, error)
	FindByLong(ctx context.Context, long string) (*Record, error)
}
