// internal/storage/sqlite/sqlite.go

package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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
	DROP TABLE comments;
	CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY,
		article_id INTEGER NOT NULL,
		text TEXT NOT NULL,
		score REAL);
`

	db.MustExec(schemaComments)

	//заполняем данными (в пределах одной транзакции)0
	//tx := db.MustBegin()
	//tx.MustExec("INSERT INTO articles (id, title, text) VALUES ($1,$2, $3)", 1, "Title 1", "This is article 1")
	//tx.MustExec("INSERT INTO articles (id, title, text) VALUES ($1,$2, $3)", 2, "Title 2", "This is article 2")
	//tx.MustExec("INSERT INTO articles (id, title, text) VALUES ($1,$2, $3)", 3, "Title 3", "This is article 3")
	//tx.MustExec("INSERT INTO articles (id, title, text) VALUES ($1,$2, $3)", 4, "Title 4", "This is article 4")
	//tx.MustExec("INSERT INTO articles (id, title, text) VALUES ($1,$2, $3)", 5, "Title 5", "This is article 5")
	//tx.MustExec("INSERT INTO articles (id, title, text) VALUES ($1,$2, $3)", 6, "Title 6", "This is article 6")
	//tx.MustExec("INSERT INTO articles (id, title, text) VALUES ($1,$2, $3)", 7, "Title 7", "This is article 7")
	//
	//// Именованные запросы могут использовать структуры,
	//// поэтому, если у вас имеется структура, (например person := &User{}),
	//// которую необходимо заполнить, Вы можете передать ее как &person:
	//// tx.NamedExec("INSERT INTO user (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &User{FirstName: "Jane", LastName: "Citizen", Email: "jane.citzen@example.com"})
	//tx.Commit()

	//region Заполнение данными
	tx2 := db.MustBegin()
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 1, "Comment 1-1", 2)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 1, "Comment 1-2", 2)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 1, "Comment 1-3", 2)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 1, "Comment 1-4", 3)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 1, "Comment 1-5", 3)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 1, "Comment 1-6", 3)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 1, "Comment 1-7", 3)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 1, "Comment 1-7", 4)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 2, "Comment2-1", 3)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 2, "Comment2-2", 3)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 2, "Comment2-3", 3)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 2, "Comment2-4", 4)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 2, "Comment2-5", 4)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 2, "Comment2-6", 4)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 2, "Comment2-7", 4)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 3, "Comment3-1", 4)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 3, "Comment3-2", 4)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 3, "Comment3-3", 4)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 3, "Comment3-4", 4)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 3, "Comment3-5", 5)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 3, "Comment3-6", 5)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 3, "Comment3-7", 5)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 4, "Comment4-1", 6)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 4, "Comment4-2", 6)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 4, "Comment4-3", 6)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 4, "Comment4-4", 7)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 4, "Comment4-5", 7)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 4, "Comment4-6", 7)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 4, "Comment4-7", 7)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 5, "Comment5-1", 7)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 5, "Comment5-2", 7)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 5, "Comment5-3", 7)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 5, "Comment5-4", 8)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 5, "Comment5-5", 8)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 5, "Comment5-6", 8)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 5, "Comment5-7", 8)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 6, "Comment6-1", 7)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 6, "Comment6-2", 5)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 6, "Comment6-3", 5)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 6, "Comment6-4", 8)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 6, "Comment6-5", 8)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 6, "Comment6-6", 8)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 6, "Comment6-7", 8)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 7, "Comment7-1", 9)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 7, "Comment7-2", 5)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 7, "Comment7-3", 5)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 7, "Comment7-4", 8)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 7, "Comment7-5", 8)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 7, "Comment7-6", 9)
	tx2.MustExec("INSERT INTO comments (article_id, text, score) VALUES ($1,$2, $3)", 7, "Comment7-7", 9)
	tx2.Commit()

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
	//min, err := s.getMinArticleId()
	//if err != nil {
	//	return nil, fmt.Errorf("%s: get min article: %w", op, err)
	//}
	//max, err := s.getMaxArticleId()
	//if err != nil {
	//	return nil, fmt.Errorf("%s: get max article: %w", op, err)
	//}

	//собственно случайное значение
	//	v := rand.Intn(max-min) + min // range is min to max
	v := 2

	if v <= 0 {
		return nil, fmt.Errorf("%s: there is no data to display (min==max)", op)
	}

	isFound := false
	var result []models.ArticleInfo

	// TODO сперва поищем в кеше redis
	//	myResult, err := s.cache.GetCachedArticle("")
	myResult, err := s.cache.GetCachedArticleAsString(strconv.Itoa(v))
	if err != nil {
		//return nil, err
		fmt.Println(time.Now(), err)
	} else {
		fmt.Println("cache found")
	}
	fmt.Println(myResult)

	counter := 0
	for !isFound || counter > 100 {
		if err := s.db.Select(&result, "SELECT id, title, text, (SELECT AVG(score) FROM comments WHERE score IS NOT NULL AND article_id= $1) as rating FROM articles WHERE id= $1", v); err != nil {
			return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
		}

		//если ничего не найдено, увеличиваем ид, и так 100 раз, потом выходим
		if len(result) > 0 {
			isFound = true
		} else {
			counter++
			v++
		}
	}

	//fmt.Printf("%+v\n", result)
	// Пишем найденное в кэш
	if myResult == "" {
		fmt.Println("set value to cache")
		s.cache.SetCachedArticle(strconv.Itoa(v), result)
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
