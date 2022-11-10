package main

import (
	"fmt"
	"log"
	"os"

	"shortener/domains/cache"
	cacheStorage "shortener/domains/cache/storage"
	"shortener/domains/generator"
	"shortener/domains/shortener"
	"shortener/domains/shortener/storage"
	"shortener/servers"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" || env == "development" {
		if err := godotenv.Load(".envs/.env"); err != nil {
			panic(err)
		}
	}

	db, err := setupDB()
	if err != nil {
		panic(err)
	}

	generatorService := generator.New()
	shortenerStore, err := storage.New(db, generatorService)
	if err != nil {
		panic(err)
	}
	shortService := shortener.New(shortenerStore)

	chStore := cacheStorage.New(setupRedisClient())
	cacheSvc := cache.New(chStore)

	server, err := servers.New(os.Getenv("BASE_URL"), db, generatorService, shortService, cacheSvc)
	if err != nil {
		panic(err)
	}

	if err := server.Router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func setupDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	fmt.Printf("\nDSN: %s\n\n", dsn)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func setupRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}
