package main

import (
	"context"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	mainS := storage.NewStorage(128)
	// connString := `postgres://user:password@localhost:5432/urlshortener`
	connString := `postgres://user:password@127.0.0.1:5433/urlshortener?sslmode=disable`
	// connString := "host=localhost port=5432 user=user password=password dbname=urlshortener sslmode=disable"
	mainRepo, err := repository.NewDB(connString)
	if err != nil {
		panic(err)
	}

	rl := ratelimit.NewRateLimiter(1000)

	mainHandler := handler.NewHandler(mainS, mainRepo)
	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", handler.RateLimit(rl, mainHandler.Shorten))
	mux.HandleFunc("/", handler.RateLimit(rl, mainHandler.GetURL))
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// http.ListenAndServe(":8080", mux)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Ждём сигнала
	<-quit
	fmt.Println("Shutting down...")

	// Graceful shutdown с таймаутом 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
	server.Shutdown(ctx)
	mainRepo.Close()
}
