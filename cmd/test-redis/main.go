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
	//region //Мой сервер для тестов с использованием стандартного http
	//http.HandleFunc("/", handler)
	//http.HandleFunc("/articles", handlerArticles)
	//http.HandleFunc("/trending", handlerTrending)
	//http.ListenAndServe(":8500", nil)
	//endregion

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

	//region Создаем объект Storage
	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	} else {
		log.Info("storage created")
	}
	//endregion

	//region Создаем http-сервер

	//region Создаем роутер
	router := chi.NewRouter()

	// Настраиваем CORS (предварительно скачиваем пакет: go get github.com/go-chi/cors)
	// Дополнительная ссылка: https://developer.github.com/v3/#cross-origin-resource-sharing
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"}, // пока что разрешаем все
		//AllowedOrigins: []string{"http://localhost:3000"}, // Use this to allow specific origin hosts
		//AllowedOrigins: []string{"https://87.242.85.156:*", "http://87.242.85.156:*"},
		//AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		//ExposedHeaders:   []string{"Link"},
		//AllowCredentials: false,
		//MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	//--------------------------------------
	//По умолчанию middleware.Logger использует свой собственный внутренний логгер,
	//который желательно переопределить, чтобы использовался наш,
	//иначе могут возникнуть проблемы — например, со сбором логов.
	//Либо можно написать собственный middleware для логирования запросов. Так и сделаем
	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.Logger)    // Логирование всех запросов. Желательно написать собственный
	router.Use(mwLogger.New(log))    // Собственный middleware для логирования запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	// РАЗОБРАТЬСЯ!!!!
	//--------------------------------------------------------------------------------
	//router.Post("/", save.New(log, storage))

	// Все пути этого роутера будут начинаться с префикса `/api`
	router.Route("/api", func(r chi.Router) {
		// Подключаем базовую аутентификацию
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			// Передаем в middleware креды
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
			// Если у вас более одного пользователя,
			// то можете добавить остальные пары по аналогии.
		}))

		//	r.Post("/", save.New(log, storage))
		//r.Post("/", save.New(log, storage))
	})
	log.Debug("Auth info", cfg.User, cfg.Password)

	router.Get("/article/{article_id}", article.GetArticle(log, storage))
	router.Get("/articles", article.GetRandArticles(log, storage))
	router.Get("/trending", article.GetTrending(log, storage))
	//router.Get("/trending", handlerArticles)

	//// Подключаем редирект-хендлер.
	//// Здесь формируем путь для обращения и именуем его параметр — {alias}.
	//// В хендлере можно получить этот параметр по указанному имени
	//// Это очень удобная и гибкая штука. Вы можете формировать и более сложные пути, например:
	//// router.Get("/v1/{user_id}/uid", redirect.New(log, storage))
	//router.Get("/{alias}", redirect.New(log, storage))
	//
	////прикручиваем ремувер
	//router.Delete("/{alias}", remove.New(log, storage))

	//endregion

	//region ЗАПУСК и ОСТАНОВКА СЕРВЕРА
	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
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

//func handlerArticles(writer http.ResponseWriter, _ *http.Request) {
//	fmt.Println("Hello articles!")
//	writer.Write([]byte("Hi articles!"))
//}
//
//func handlerTrending(writer http.ResponseWriter, _ *http.Request) {
//	fmt.Println("Hello Trending!")
//	writer.Write([]byte("Hi Trending!"))
//}
