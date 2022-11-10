package servers

import (
	"shortener/domains/cache"
	"shortener/domains/generator"
	"shortener/domains/shortener"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Server constructs the basic server components
type Server struct {
	baseURL      string
	DB           *gorm.DB
	Redis        *redis.Client
	Router       *gin.Engine
	genService   generator.IGenerateService
	shortSvc     shortener.IShortenerService
	cacheService cache.ICacheService
}

// New creates a new server instance injected with necessary service
func New(baseURL string,
	db *gorm.DB,
	genService generator.IGenerateService,
	shortSvc shortener.IShortenerService,
	cacheSvc cache.ICacheService) (*Server, error) {
	s := &Server{
		baseURL:      baseURL,
		Router:       gin.Default(),
		genService:   genService,
		shortSvc:     shortSvc,
		cacheService: cacheSvc,
	}

	s.Router.POST("/", s.shorten)
	s.Router.GET("/:short", s.visit)
	return s, nil
}
