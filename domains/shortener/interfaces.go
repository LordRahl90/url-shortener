package shortener

import (
	"context"

	"shortener/domains/entities"
)

// IShortenerService interface for shortener service
type IShortenerService interface {
	Create(ctx context.Context, record *entities.Shortener) error
	FindByShort(ctx context.Context, short string) (*entities.Shortener, error)
	FindByLong(ctx context.Context, long string) (*entities.Shortener, error)
}
