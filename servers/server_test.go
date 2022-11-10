package servers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"shortener/domains/cache"
	cacheStore "shortener/domains/cache/storage"
	"shortener/domains/entities"
	"shortener/domains/generator"
	"shortener/domains/shortener"
	"shortener/domains/shortener/storage"
	"shortener/requests"
	"shortener/responses"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	baseURL = "https://short.io"

	db         *gorm.DB
	genService generator.IGenerateService
	shortSvc   shortener.IShortenerService

	cacheSvc cache.ICacheService

	server *Server
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		cleanup()
		os.Exit(code)
	}()
	d, err := setupTestDB()
	if err != nil {
		panic(err)
	}
	db = d
	genService = generator.New()
	store, err := storage.New(db, genService)
	if err != nil {
		panic(err)
	}
	shortSvc = shortener.New(store)

	chStore := cacheStore.New(setupRedisClient())
	cacheSvc = cache.New(chStore)

	s, err := New(baseURL, db, genService, shortSvc, cacheSvc)
	if err != nil {
		panic(err)
	}
	server = s

	code = m.Run()
}

func TestShortenLink(t *testing.T) {
	req := &requests.Shortener{
		Link: gofakeit.URL(),
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotNil(t, b)

	res := handleRequest(t, http.MethodPost, "/", b)
	require.Equal(t, http.StatusCreated, res.Code)

	var response *responses.Shortener
	err = json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.Equal(t, response.Link, req.Link)
	spl := strings.Split(response.Short, baseURL)
	require.Len(t, spl, 2)
	short := spl[1][1:]
	assert.True(t, len(short) <= 10) //remove the backslash

	// make sure the data is truly registered
	data, err := shortSvc.FindByShort(context.Background(), short)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	res = handleRequest(t, http.MethodGet, "/"+short, nil)
	require.Equal(t, http.StatusMovedPermanently, res.Code)
}

func TestShortenWithInvalidJSON(t *testing.T) {
	b := []byte(`
	{"link":"https://www.corporatee-tailers.info/exploit",}
	`)
	res := handleRequest(t, http.MethodPost, "/", b)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestShortenWithBAdURL(t *testing.T) {
	b := []byte(`{"link":"corporatee-tailers hello exploit"}`)
	res := handleRequest(t, http.MethodPost, "/", b)
	require.Equal(t, http.StatusBadRequest, res.Code)
	exp := `{"error":"parse \"corporatee-tailers hello exploit\": invalid URI for request","success":false}`
	require.Equal(t, exp, res.Body.String())
}

func TestVisitShortLink(t *testing.T) {
	ctx := context.Background()
	svcRec := &entities.Shortener{
		LongText: gofakeit.URL(),
	}
	err := shortSvc.Create(ctx, svcRec)
	require.NoError(t, err)
	require.True(t, len(svcRec.ShortText) <= 10)

	res := handleRequest(t, http.MethodGet, "/"+svcRec.ShortText, nil)
	require.Equal(t, http.StatusMovedPermanently, res.Code)
}

func TestVisitNonExistentShortLink(t *testing.T) {
	res := handleRequest(t, http.MethodGet, "/"+uuid.NewString(), nil)
	require.Equal(t, http.StatusNotFound, res.Code)
}

func handleRequest(t *testing.T, method, path string, payload []byte) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	var (
		req *http.Request
		err error
	)
	if len(payload) > 0 {
		req, err = http.NewRequest(method, path, bytes.NewBuffer(payload))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	require.NoError(t, err)
	server.Router.ServeHTTP(w, req)

	return w
}

func setupTestDB() (*gorm.DB, error) {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:@tcp(127.0.0.1:3306)/shortener?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "test_user:password@tcp(127.0.0.1:33306)/shortener?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

func cleanup() {
	db.Exec("DELETE FROM records")
}
