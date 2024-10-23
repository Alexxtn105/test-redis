package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"test-redis/internal/lib/logger/sl"
	"test-redis/internal/storage"
)

// URLGetter is an interface for getting url by alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func GetArticle(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.article.GteArticle"

		//пишем в лог
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Роутер chi позволяет делать вот такие финты -
		// получать GET-параметры по их именам.
		// Имена определяются при добавлении хэндлера в роутер, это будет ниже.
		article_id := chi.URLParam(r, "article_id")
		if article_id == "" {
			log.Info("article_id is empty")

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		// Находим URL по алиасу в БД
		resURL, err := urlGetter.GetURL(article_id)
		if errors.Is(err, storage.ErrURLNotFound) {
			// Не нашли URL, сообщаем об этом клиенту
			log.Info("url not found", "article_id", article_id)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			// Не удалось осуществить поиск
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		// Делаем редирект на найденный URL
		http.Redirect(w, r, resURL, http.StatusFound)
		// В последней строчке делаем редирект со статусом http.StatusFound — код HTTP 302. Он обычно используется для временных перенаправлений, а не постоянных, за которые отвечает 301.
		// Наш сервис может перенаправлять на разные URL в зависимости от ситуации
		// (мы ведь можем удалить или изменить сохраненный URL),
		// поэтому есть смысл использовать именно http.StatusFound.
		// Это важно для систем кэширования и поисковых машин —
		// они обычно кэшируют редиректы с кодом 301, то есть считают их постоянными.
		// Нам такое поведение не нужно.
	}
}
