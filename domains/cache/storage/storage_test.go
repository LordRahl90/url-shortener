package storage

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

var (
	client *redis.Client
	store  ICacheStore
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		client.Close()
		os.Exit(code)
	}()
	client = setupRedisClient()
	if client == nil {
		panic("redis not connected")
	}
	store = New(client)
	code = m.Run()
}

func TestSaveRecord(t *testing.T) {
	ctx := context.Background()
	key := "12345"
	value := gofakeit.URL()

	err := store.Save(ctx, key, value)
	require.NoError(t, err)

	v, err := store.Find(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, v)

}

func setupRedisClient() *redis.Client {
	env := os.Getenv("ENVIRONMENT")
	if env == "cicd" {
		return redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	}
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "password123",
		DB:       0,
	})
}
