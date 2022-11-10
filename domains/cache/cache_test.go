package cache

import (
	"context"
	"fmt"
	"os"
	"testing"

	"shortener/domains/cache/storage"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

var (
	store storage.ICacheStore

	mockDB = make(map[string]string)
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()
	store = &MockCacheStore{
		SaveFunc: func(ctx context.Context, key, value string) error {
			mockDB[key] = value
			return nil
		},
		FindFunc: func(ctx context.Context, key string) (string, error) {
			return mockDB[key], nil
		},
	}
	code = m.Run()
}

func TestSaveNewRecord(t *testing.T) {
	ctx := context.Background()
	svc := New(store)
	key := "12345"
	value := gofakeit.URL()

	err := svc.Save(ctx, key, value)
	require.NoError(t, err)

	v, err := svc.Find(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, v)
}

func TestSaveWithError(t *testing.T) {
	ctx := context.Background()

	st := &MockCacheStore{
		SaveFunc: func(ctx context.Context, key, value string) error {
			return fmt.Errorf("cannot save record")
		},
	}

	svc := New(st)
	key := "12345"
	value := gofakeit.URL()

	require.EqualError(t, svc.Save(ctx, key, value), "cannot save record")
}

func TestSaveWithMockNotInitialized(t *testing.T) {
	ctx := context.Background()

	st := &MockCacheStore{}

	svc := New(st)
	key := "12345"
	value := gofakeit.URL()

	require.EqualError(t, svc.Save(ctx, key, value), errMockNotInitialized.Error())
}

func TestFindWithMockNotInitialized(t *testing.T) {
	ctx := context.Background()

	st := &MockCacheStore{}

	svc := New(st)
	key := "12345"

	res, err := svc.Find(ctx, key)

	require.EqualError(t, err, errMockNotInitialized.Error())
	require.Empty(t, res)
}

func TestFindWithError(t *testing.T) {
	ctx := context.Background()

	st := &MockCacheStore{
		FindFunc: func(ctx context.Context, key string) (string, error) {
			return "", fmt.Errorf("cannot get this value")
		},
	}

	svc := New(st)
	key := "12345"
	res, err := svc.Find(ctx, key)

	require.EqualError(t, err, "cannot get this value")
	require.Empty(t, res)
}
