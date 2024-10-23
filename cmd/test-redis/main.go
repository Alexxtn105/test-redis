package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"test-redis/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//region Загружаем конфигурацию
	cfg := config.MustLoadFetchFlag() // ...или с использованием параметра командной строки
	fmt.Println("Конфигурация загружена успешно")
	//endregion

	//region Создаем логгер
	log := setupLogger(cfg.Env)
	//добавим параметр env с помощью метода log.With
	log = log.With(slog.String("env", cfg.Env))                          // к каждому сообщению будет добавляться поле с информацией о текущем окружении
	log.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")
	//endregion

	//http.HandleFunc("/", handler)
	http.HandleFunc("/articles", handlerArticles)
	http.HandleFunc("/trending", handlerTrending)
	http.ListenAndServe(":8500", nil)
}

// setupLogger создает логгер в зависимости от окружения с разными параметрами — TextHandler / JSONHandler и уровень LevelDebug / LevelInfo
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return log
}

func handler(writer http.ResponseWriter, _ *http.Request) {
	fmt.Println("Hello from Redis!")
	writer.Write([]byte("Hi there!"))
}

func handlerArticles(writer http.ResponseWriter, _ *http.Request) {
	fmt.Println("Hello articles!")
	writer.Write([]byte("Hi articles!"))
}
func handlerTrending(writer http.ResponseWriter, _ *http.Request) {
	fmt.Println("Hello Trending!")
	writer.Write([]byte("Hi Trending!"))
}
