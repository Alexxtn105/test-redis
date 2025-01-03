package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-redis/internal/cache/redisCache"
	"test-redis/internal/config"
	"test-redis/internal/lib/logger/sl"
	"test-redis/internal/storage/sqlite"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	mwLogger "test-redis/internal/http-server/middleware/logger"

	"test-redis/internal/http-server/handlers"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//region // Сервер для тестов с использованием стандартного http
	//http.HandleFunc("/", handler)
	//http.HandleFunc("/articles", handlerArticles)
	//http.HandleFunc("/trending", handlerTrending)
	//http.ListenAndServe(":8500", nil)
	//endregion

	//region Загружаем конфигурацию
	cfg := config.MustLoadFetchFlag() // ...или с использованием параметра командной строки
	fmt.Println("time=", time.Now(), "Конфигурация загружена успешно")
	//endregion

	//region Создаем логгер
	fmt.Println("time=", time.Now(), "Создание логгера")
	log := setupLogger(cfg.Env)
	//добавим параметр env с помощью метода log.With
	log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении
	log.Debug("logger debug mode enabled")
	//endregion

	//region Создаем объект кэша Redis
	log.Info("initializing cache", slog.String("address", cfg.Cache.Address)) // Помимо сообщения выведем параметр с адресом
	cache, err := redisCache.NewCache(cfg.Cache.Address, cfg.Cache.Password, cfg.Cache.DB)
	if err != nil {
		log.Error("failed to initialize cache", sl.Err(err))
	} else {
		log.Info("cache created")
	}
	//endregion

	//region Создаем объект Storage Sqlite 3
	log.Info("initializing storage", slog.String("storage_path", cfg.StoragePath)) // Помимо сообщения выведем параметр с адресом
	storage, err := sqlite.NewStorage(cfg.StoragePath, cache)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	} else {
		log.Info("storage created")
	}
	//endregion

	//region Создаем роутер
	router := chi.NewRouter()

	// Настраиваем CORS (предварительно скачиваем пакет: go get github.com/go-chi/cors)
	// Дополнительная ссылка: https://developer.github.com/v3/#cross-origin-resource-sharing
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"}, // пока что разрешаем все
		//AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		//ExposedHeaders:   []string{"Link"},
		//AllowCredentials: false,
		//MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// По умолчанию middleware.Logger использует свой собственный внутренний логгер,
	// который желательно переопределить, чтобы использовался наш,
	// иначе могут возникнуть проблемы — например, со сбором логов.
	// Либо можно написать собственный middleware для логирования запросов. Так и сделаем
	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.Logger)    // Логирование всех запросов. Желательно написать собственный
	router.Use(mwLogger.New(log))    // Собственный middleware для логирования запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер url поступающих запросов

	// Прописываем маршруты с параметром {article_id}.
	// В хендлере можно получить этот параметр по указанному имени
	// Это очень удобная и гибкая штука. Можно формировать и более сложные пути, например:
	// router.Get("/v1/{user_id}/uid", redirect.New(log, storage))

	router.Get("/article/{article_id}", article.GetArticle(log, storage))
	router.Get("/articles", article.GetRandArticles(log, storage))
	//router.Get("/articles", article.GetTestData(log))
	router.Get("/test", article.GetTestData(log))
	router.Get("/users/{user_id}", article.GetUserById(log))
	//endregion

	//region ЗАПУСК и ОСТАНОВКА СЕРВЕРА
	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				log.Error("failed to start server", slog.String("error", err.Error()))
			}
		}
	}()
	log.Info("server started")

	// ждем, пока в канал не придет сигнал с остановкой сервера
	<-done
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))
		return
	}

	// TODO: close storage
	//...

	log.Info("server stopped")
	//endregion

}

// setupLogger Создает логгер в зависимости от окружения с разными параметрами — TextHandler / JSONHandler и уровень LevelDebug / LevelInfo
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
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
