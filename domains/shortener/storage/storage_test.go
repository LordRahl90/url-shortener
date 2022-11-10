package storage

import (
	"context"
	"os"
	"testing"

	"shortener/domains/generator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	dbErr error
	store IShortenerStorage
	gen   = generator.New()
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		cleanup()
		os.Exit(code)
	}()

	db, dbErr = setupTestDB()
	if dbErr != nil {
		panic(dbErr)
	}
	store, dbErr = New(db, gen)
	if dbErr != nil {
		panic(dbErr)
	}
	code = m.Run()
}

func TestCreateShortener(t *testing.T) {
	ctx := context.Background()
	r := &Record{
		LongText: "https://google.com?q=phones",
	}
	t.Cleanup(func() {
		db.Exec("DELETE FROM records WHERE id = ?", r.ID)
	})

	err := store.Create(ctx, r)
	require.NoError(t, err)
	assert.Len(t, r.ShortText, 10)

	r1 := &Record{
		LongText: "https://google.com?q=phones",
	}

	err = store.Create(ctx, r1)
	require.NoError(t, err)

	assert.Equal(t, r1.ID, r.ID)
}

func TestFindRecordByLongURLNotFound(t *testing.T) {
	ctx := context.Background()
	res, err := store.FindByLong(ctx, "https://google.com?q=phones+testing")
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func TestFindRecordByLongURL(t *testing.T) {
	ctx := context.Background()
	r := &Record{
		LongText: "https://google.com?q=phones",
	}
	t.Cleanup(func() {
		db.Exec("DELETE FROM records WHERE id = ?", r.ID)
	})

	err := store.Create(ctx, r)
	require.NoError(t, err)
	assert.Len(t, r.ShortText, 10)

	res, err := store.FindByLong(ctx, r.LongText)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	assert.Equal(t, r.ID, res.ID)
	assert.Equal(t, r.ShortText, res.ShortText)
}

func TestFindRecordByShortURL(t *testing.T) {
	ctx := context.Background()
	r := &Record{
		LongText: "https://google.com?q=phones",
	}
	t.Cleanup(func() {
		db.Exec("DELETE FROM records WHERE id = ?", r.ID)
	})

	err := store.Create(ctx, r)
	require.NoError(t, err)
	assert.Len(t, r.ShortText, 10)

	res, err := store.FindByShort(ctx, r.ShortText)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	assert.Equal(t, r.ID, res.ID)
	assert.Equal(t, r.ShortText, res.ShortText)
}

func TestFindRecordByShortNotExist(t *testing.T) {
	ctx := context.Background()
	res, err := store.FindByShort(ctx, "12HrElrRaK")
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func setupTestDB() (*gorm.DB, error) {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:@tcp(127.0.0.1:3306)/shortener?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "test_user:password@tcp(127.0.0.1:33306)/shortener?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func cleanup() {
	db.Exec("DELETE FROM records")
}
