package servers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"shortener/domains/generator"
	"shortener/domains/shortener"
	"shortener/domains/shortener/storage"
	"shortener/requests"
	"shortener/responses"

	"github.com/brianvoe/gofakeit/v6"
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

	server *Server
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		// cleanup()
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

	s, err := New(baseURL, db, genService, shortSvc)
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

func cleanup() {
	db.Exec("DELETE FROM records")
}
