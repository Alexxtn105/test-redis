// internal/storage/sqlite/sqlite.go

package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"strconv"
	"test-redis/internal/cache/redisCache"
	"test-redis/internal/models"
	"test-redis/internal/storage"
	"time"
)

// Storage Структура объекта Storage
type Storage struct {
	//	db *sql.DB //из пакета "database/sql"
	db    *sqlx.DB //из пакета "database/sql"
	cache *redisCache.Cache
}

// NewStorage Конструктор объекта Storage
func NewStorage(storagePath string, cache *redisCache.Cache) (*Storage, error) {
	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок

	// Подключаемся к БД (сделал с использованием sqlx - https://github.com/joncrlsn/go-examples/blob/master/sqlx-sqlite.go)
	db, err := sqlx.Connect("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	} else {
		//fmt.Println("db connected")
	}

	// TODO: можно прикрутить миграции, для тренировки
	// создаем таблицу, если ее еще нет
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers0;  pq will exec them all, sqlite3 won't, ymmv
	var schemaArticles = `
	--таблица статей
	--DROP TABLE articles;
	CREATE TABLE IF NOT EXISTS articles(
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL UNIQUE,
		text TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_theme ON articles(title);
`
	db.MustExec(schemaArticles)

	var schemaComments = `
	--таблица комментариев
	--DROP TABLE comments;
	CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY,
		article_id INTEGER NOT NULL,
		text TEXT NOT NULL,
		score REAL);
`
	db.MustExec(schemaComments)

	//region Заполнение данными статей
	//
	//tx1 := db.MustBegin()
	//title := ""
	//article := ""
	//for i := 1; i < 100; i++ {
	//	title = fmt.Sprintf("Title %d", i)
	//	article = fmt.Sprintf("This is article %d", i)
	//	tx1.MustExec("INSERT INTO articles (title, text) VALUES ($1,$2)", title, article)
	//}
	//tx1.Commit()
	//endregion
	//region Заполнение данными комментариев
	//min := 0
	//max := 100
	//
	//tx2 := db.MustBegin()
	//
	//comment := ""
	//for i := 1; i < 100; i++ {
	//	for j := 1; j < 100; j++ {
	//		v := rand.Intn(max-min) + min // range is min to max
	//
	//		comment = fmt.Sprintf("comment %d-%d", i, j)
	//		tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", i, comment, v)
	//	}
	//}
	//tx2.Commit()

	//endregion

	//// Тестовый запрос в базу данных, результаты сохраним в слайс []models.ArticleInfo
	//var articles []models.ArticleInfo
	//db.Select(&articles, "SELECT * FROM articles ORDER BY title ASC")
	////выводим первые две
	//article1, article2 := articles[0], articles[1]
	//fmt.Printf("Article 1: %#v\nArticle 2: %#v\n", article1, article2)

	return &Storage{db: db, cache: cache}, nil
}

// getMinArticleId Получить минимальный ИД из таблицы статей
func (s *Storage) getMinArticleId() (int, error) {
	const op = "storage.sqlite.getMinArticleId"

	var cnt []int
	if err := s.db.Select(&cnt, "SELECT MIN(id) FROM articles"); err != nil {
		return 0, fmt.Errorf("%s: select query: %w", op, err)
	}
	if len(cnt) == 0 {
		return 0, nil
	}

	return cnt[0], nil
}

// getMaxArticleId Получить максимальный ИД из таблицы статей
func (s *Storage) getMaxArticleId() (int, error) {
	const op = "storage.sqlite.getMaxArticleId"

	var cnt []int
	if err := s.db.Select(&cnt, "SELECT MAX(id) FROM articles"); err != nil {
		return 0, fmt.Errorf("%s: select query: %w", op, err)
	}
	if len(cnt) == 0 {
		return 0, nil
	}
	return cnt[0], nil
}

// GetRandomData Получить случайную статью из таблицы
func (s *Storage) GetRandomData() ([]models.ArticleInfo, error) {
	const op = "storage.sqlite.GetRandomData"

	//берем случайное число в диапазоне от минимального до максимального ид статьи
	min := 1
	max := 100
	//min, err := s.getMinArticleId()
	//if err != nil {
	//	return nil, fmt.Errorf("%s: get min article: %w", op, err)
	//}
	//max, err := s.getMaxArticleId()
	//if err != nil {
	//	return nil, fmt.Errorf("%s: get max article: %w", op, err)
	//}

	//собственно случайное значение
	v := rand.Intn(max-min) + min // range is min to max
	//v := 2

	if v <= 0 {
		return nil, fmt.Errorf("%s: there is no data to display (min==max)", op)
	}

	var result []models.ArticleInfo

	// Сперва поищем в кеше redis
	result, err := s.cache.GetCachedArticle(strconv.Itoa(v))
	if err != nil {
		fmt.Println(time.Now(), err)
	}

	// Если ничего не найдено, увеличиваем ид, и так 100 раз, потом выходим
	if result == nil {
		isFound := false
		counter := 0
		for !isFound || counter > 100 {
			if err := s.db.Select(&result, "SELECT id, title, text, (SELECT AVG(score) FROM comments WHERE score IS NOT NULL AND article_id= $1) as rating FROM articles WHERE id= $1", v); err != nil {
				return nil, fmt.Errorf("%s: select: %w", op, err)
			}

			if len(result) > 0 {
				isFound = true

				// Пишем найденное в кэш
				fmt.Println("set value to cache")

				// преобразуем структуру в []byte для хранения в кэше
				res, err := json.Marshal(result)
				if err != nil {
					// TODO - разобраться с ошибкой
					return result, fmt.Errorf("%s: cannot set value to cache: %w", op, err)
				} else {
					s.cache.SetCachedArticle(strconv.Itoa(v), res)
				}
			} else {
				counter++
				v++
			}
		}
	}

	return result, nil
}

// GetData - получить данные
func (s *Storage) GetData(id string) (string, error) {
	const op = "storage.sqlite.GetData"

	// Подготавливаем запрос (проверка корректности синтаксиса)
	stmt, err := s.db.Prepare("SELECT text FROM articles WHERE id = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var result string

	//в параметрах используем указатель, чтобы получить результаты
	err = stmt.QueryRow(id).Scan(&result)

	//если строки не найдено - возвращаем пустую строку
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrDataNotFound
	}

	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return result, nil
}
