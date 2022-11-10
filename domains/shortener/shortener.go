package shortener

import (
	"context"
	"fmt"
	"net/url"

	"shortener/domains/entities"
	"shortener/domains/shortener/storage"
)

var (
	_ IShortenerService = (*ShortenerService)(nil)
)

// ShortenerService concrete implementation of the IShortenerService
type ShortenerService struct {
	store storage.IShortenerStorage
}

// New returns a new instance of shortener storage
func New(store storage.IShortenerStorage) IShortenerService {
	return &ShortenerService{
		store: store,
	}
}

// Create implements IShortenerService
func (s *ShortenerService) Create(ctx context.Context, record *entities.Shortener) error {
	// make sure we are getting a valid url
	fmt.Printf("\n\nLong: %s\n\n", record.LongText)
	link, err := url.ParseRequestURI(record.LongText)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nScheme: %s\n\n", link)

	dbRecord := record.ToDBEntity()
	if err := s.store.Create(ctx, dbRecord); err != nil {
		return err
	}
	record.ShortText = dbRecord.ShortText
	return nil
}

// FindByLong fetches the record with the long text
func (s *ShortenerService) FindByLong(ctx context.Context, long string) (*entities.Shortener, error) {
	res, err := s.store.FindByLong(ctx, long)
	if err != nil {
		return nil, err
	}
	return entities.ShortenerFromDBEntity(res), nil
}

// FindByShort fetches a record with the short text
func (s *ShortenerService) FindByShort(ctx context.Context, short string) (*entities.Shortener, error) {
	res, err := s.store.FindByShort(ctx, short)
	if err != nil {
		return nil, err
	}
	return entities.ShortenerFromDBEntity(res), nil
}
