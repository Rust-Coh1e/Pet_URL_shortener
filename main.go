package main

import (
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "url-shortener/docs"
	"url-shortener/handler"
	"url-shortener/ratelimit"
	"url-shortener/repository"
	"url-shortener/storage"
)

// @title          URL Shortener API
// @version        1.0
// @description    Сервис сокращения ссылок
// @host           localhost:8080
func main() {
	fmt.Println("Starting...")

	mainS := storage.NewStorage()
	// connString := `postgres://user:password@localhost:5432/urlshortener`
	connString := `postgres://user:password@127.0.0.1:5433/urlshortener?sslmode=disable`
	// connString := "host=localhost port=5432 user=user password=password dbname=urlshortener sslmode=disable"
	mainRepo, err := repository.NewDB(connString)
	if err != nil {
		panic(err)
	}

	rl := ratelimit.NewRateLimiter(100)

	mainHandler := handler.NewHandler(mainS, mainRepo)
	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", handler.RateLimit(rl, mainHandler.Shorten))
	mux.HandleFunc("/", handler.RateLimit(rl, mainHandler.GetURL))
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.ListenAndServe(":8080", mux)
}
