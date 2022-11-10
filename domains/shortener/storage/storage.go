package storage

import (
	"context"
	"errors"
	"sync"
	"time"

	"shortener/domains/generator"

	"gorm.io/gorm"
)

var _ IShortenerStorage = (*Storage)(nil)

// Storage repository to manage shortened url records
type Storage struct {
	db              *gorm.DB
	mx              sync.Mutex
	generateService generator.IGenerateService
}

// New returns a new instance of IStorage
// A machine ID is needed here to prevent clashes across different instances.
func New(db *gorm.DB, generateService generator.IGenerateService) (IShortenerStorage, error) {
	if err := db.AutoMigrate(&Record{}); err != nil {
		return nil, err
	}
	return &Storage{
		db:              db,
		mx:              sync.Mutex{},
		generateService: generateService,
	}, nil
}

// Create implements creating a new record
// we need to get a mutex lock so we can avoid time collissions.
// this could also be accomplished by having a separate package to generate IDs
func (s *Storage) Create(ctx context.Context, r *Record) error {
	res, err := s.FindByLong(ctx, r.LongText)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if res.ID > 0 {
		r.ID = res.ID
		r.ShortText = res.ShortText
		return nil
	}
	r.ID = uint64(time.Now().Nanosecond())
	r.ShortText = s.generateService.Generate(ctx, int(r.ID))

	return s.db.WithContext(ctx).Create(&r).Error
}

// FindByLong finds a record by the long string
func (s *Storage) FindByLong(ctx context.Context, long string) (result *Record, err error) {
	err = s.db.WithContext(ctx).Where("long_text = ?", long).First(&result).Error
	return
}

// FindByShort finds a record by the short string
func (s *Storage) FindByShort(ctx context.Context, short string) (result *Record, err error) {
	err = s.db.WithContext(ctx).Where("short_text = ?", short).First(&result).Error
	return
}
