//internal/http-server/handlers/article.go

package article

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io/ioutil"
	"log/slog"
	"net/http"
	resp "test-redis/internal/lib/api/response"
	"test-redis/internal/lib/logger/sl"
	"test-redis/internal/models"
	"test-redis/internal/storage"
)

// DataGetter is an interface for getting data by Id.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=DataGetter
type DataGetter interface {
	GetData(id string) (string, error)
	GetRandomData() ([]models.ArticleInfo, error)
}

// GetRandArticles Получить статьи по их ид
func GetRandArticles(log *slog.Logger, dataGetter DataGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.article.GetArticles"

		//пишем в лог
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Находим статью в БД
		resData, err := dataGetter.GetRandomData()
		if errors.Is(err, storage.ErrDataNotFound) {
			// Не нашли, сообщаем об этом клиенту
			log.Info("data not found")
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			// Не удалось осуществить поиск
			log.Error("failed to get data", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		//		log.Info("got data", slog.String("data", resData))
		log.Info("got data")

		//пишем в ответ
		res, err := json.Marshal(resData)
		w.Write(res)

		// Делаем редирект на найденный URL
		//http.Redirect(w, r, resData, http.StatusFound)

		// В последней строчке делаем редирект со статусом http.StatusFound — код HTTP 302. Он обычно используется для временных перенаправлений, а не постоянных, за которые отвечает 301.
		// Наш сервис может перенаправлять на разные URL в зависимости от ситуации
		// (мы ведь можем удалить или изменить сохраненный URL),
		// поэтому есть смысл использовать именно http.StatusFound.
		// Это важно для систем кэширования и поисковых машин —
		// они обычно кэшируют редиректы с кодом 301, то есть считают их постоянными.
		// Нам такое поведение не нужно.
	}
}

// GetArticle Получить статью по ее ид
func GetArticle(log *slog.Logger, dataGetter DataGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.article.GetArticle"

		//пишем в лог
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Роутер chi позволяет делать вот такие финты - получать GET-параметры по их именам.
		// Имена определяются при добавлении хэндлера в роутер.
		articleId := chi.URLParam(r, "article_id")
		if articleId == "" {
			log.Info("article_id is empty")
			render.JSON(w, r, resp.Error("not found"))
			return
		}

		// Находим статью в БД
		resData, err := dataGetter.GetData(articleId)
		if errors.Is(err, storage.ErrDataNotFound) {
			// Не нашли, сообщаем об этом клиенту
			log.Info("data not found", "article_id", articleId)
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			// Не удалось осуществить поиск
			log.Error("failed to get data", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("got data", slog.String("data", resData))

		//пишем в ответ
		w.Write([]byte("response: " + resData))

		// Делаем редирект на найденный URL
		//http.Redirect(w, r, resData, http.StatusFound)

		// В последней строчке делаем редирект со статусом http.StatusFound — код HTTP 302. Он обычно используется для временных перенаправлений, а не постоянных, за которые отвечает 301.
		// Наш сервис может перенаправлять на разные URL в зависимости от ситуации
		// (мы ведь можем удалить или изменить сохраненный URL),
		// поэтому есть смысл использовать именно http.StatusFound.
		// Это важно для систем кэширования и поисковых машин —
		// они обычно кэшируют редиректы с кодом 301, то есть считают их постоянными.
		// Нам такое поведение не нужно.
	}
}

func GetTestData(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := "https://jsonplaceholder.typicode.com/users"
		method := "GET"

		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			log.Error("failed to create request", sl.Err(err))
			return
		}
		res, err := client.Do(req)
		if err != nil {
			log.Error("failed to get data", sl.Err(err))
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error("failed to read data", sl.Err(err))
			//fmt.Println(err)
			return
		}
		w.Write(body)
	}
}

//func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
//	render.JSON(w, r, Response{
//		Response: resp.OK(),
//		Alias:    alias,
//	})
//}
