package shortener

import (
	"context"
	"os"
	"testing"
	"time"

	"shortener/domains/entities"
	"shortener/domains/shortener/storage"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()

	code = m.Run()
}

func TestCreateRecord(t *testing.T) {
	store := &MockStorage{
		CreateFunc: func(ctx context.Context, r *storage.Record) error {
			r.ID = uint64(time.Nanosecond)
			r.ShortText = gofakeit.Word()
			return nil
		},
	}
	svc := New(store)
	ctx := context.Background()
	r := &entities.Shortener{
		LongText: gofakeit.ImageURL(100, 100),
	}
	err := svc.Create(ctx, r)
	require.NoError(t, err)
	require.NotEmpty(t, r.ShortText)
	assert.True(t, len(r.ShortText) < 10)
}

func TestCreateWithError(t *testing.T) {
	store := &MockStorage{
		CreateFunc: func(ctx context.Context, r *storage.Record) error {
			return gorm.ErrRecordNotFound
		},
	}
	svc := New(store)
	ctx := context.Background()
	r := &entities.Shortener{
		LongText: gofakeit.ImageURL(100, 100),
	}
	err := svc.Create(ctx, r)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func TestCreateRecordWithBadURL(t *testing.T) {
	store := &MockStorage{
		CreateFunc: func(ctx context.Context, r *storage.Record) error {
			r.ID = uint64(time.Nanosecond)
			r.ShortText = gofakeit.Word()
			return nil
		},
	}
	svc := New(store)
	ctx := context.Background()
	r := &entities.Shortener{
		LongText: gofakeit.Paragraph(2, 2, 10, " "),
	}

	err := svc.Create(ctx, r)
	require.NotNil(t, err)
}

func TestCreateWithUndefinedMock(t *testing.T) {
	store := &MockStorage{}
	svc := New(store)
	ctx := context.Background()
	r := &entities.Shortener{
		LongText: gofakeit.ImageURL(100, 100),
	}
	err := svc.Create(ctx, r)
	require.EqualError(t, err, errMockNotInitialized.Error())
}

func TestFindByShort(t *testing.T) {
	store := &MockStorage{
		FindByShortFunc: func(ctx context.Context, short string) (*storage.Record, error) {
			return &storage.Record{
				ID:        uint64(time.Nanosecond),
				LongText:  gofakeit.URL(),
				ShortText: short,
			}, nil
		},
	}
	svc := New(store)
	ctx := context.Background()

	res, err := svc.FindByShort(ctx, gofakeit.Word())
	require.NoError(t, err)
	require.NotEmpty(t, res)
	assert.NotEmpty(t, res.LongText)
}

func TestFindByShortNotFound(t *testing.T) {
	store := &MockStorage{
		FindByShortFunc: func(ctx context.Context, short string) (*storage.Record, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := New(store)
	ctx := context.Background()

	res, err := svc.FindByShort(ctx, gofakeit.Word())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func TestFindShortEmptyMock(t *testing.T) {
	store := &MockStorage{}
	svc := New(store)
	ctx := context.Background()

	res, err := svc.FindByShort(ctx, gofakeit.Word())
	require.EqualError(t, err, errMockNotInitialized.Error())
	require.Empty(t, res)
}

func TestFindByLong(t *testing.T) {
	store := &MockStorage{
		FindByLongFunc: func(ctx context.Context, long string) (*storage.Record, error) {
			return &storage.Record{
				ID:        uint64(time.Nanosecond),
				LongText:  long,
				ShortText: gofakeit.Word(),
			}, nil
		},
	}
	svc := New(store)
	ctx := context.Background()
	res, err := svc.FindByLong(ctx, gofakeit.URL())
	require.NoError(t, err)
	require.NotEmpty(t, res)
	assert.NotEmpty(t, res.ShortText)
	assert.True(t, len(res.ShortText) < 10)
}

func TestFindByLongNotFound(t *testing.T) {
	store := &MockStorage{
		FindByLongFunc: func(ctx context.Context, long string) (*storage.Record, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := New(store)
	ctx := context.Background()

	res, err := svc.FindByLong(ctx, gofakeit.URL())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func TestFindLongEmptyMock(t *testing.T) {
	store := &MockStorage{}
	svc := New(store)
	ctx := context.Background()

	res, err := svc.FindByLong(ctx, gofakeit.URL())
	require.EqualError(t, err, errMockNotInitialized.Error())
	require.Empty(t, res)
}
