// internal/storage/sqlite/sqlite.go

package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"test-redis/internal/models"
	"test-redis/internal/storage"
)

// Storage Структура объекта Storage
type Storage struct {
	//	db *sql.DB //из пакета "database/sql"
	db *sqlx.DB //из пакета "database/sql"
}

var schemaArticles = `
	--таблица статей
	CREATE TABLE IF NOT EXISTS articles(
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL UNIQUE,
		text TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_theme ON articles(title);
`

var schemaComments = `
	--таблица комментариев
	CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY,
		text TEXT NOT NULL,
		score INTEGER);
`

// NewStorage Конструктор объекта Storage
func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок

	// Подключаемся к БД (сделал с использованием sqlx - https://github.com/joncrlsn/go-examples/blob/master/sqlx-sqlite.go)
	db, err := sqlx.Connect("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	} else {
		fmt.Println("db connected")
	}

	// TODO: можно прикрутить миграции, для тренировки
	// создаем таблицу, если ее еще нет
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers0;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schemaArticles)
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

	//// Тестовый запрос в базу данных, результаты сохраним в слайс []models.ArticleInfo
	//var articles []models.ArticleInfo
	//db.Select(&articles, "SELECT * FROM articles ORDER BY title ASC")
	////выводим первые две
	//article1, article2 := articles[0], articles[1]
	//fmt.Printf("Article 1: %#v\nArticle 2: %#v\n", article1, article2)

	return &Storage{db: db}, nil
}

//func NewStorage_OLD(storagePath string) (*Storage, error) {
//	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок
//
//	// Подключаемся к БД
//	db, err := sql.Open("sqlite3", storagePath)
//
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	// TODO: можно прикрутить миграции, для тренировки
//	// создаем таблицу, если ее еще нет
//	stmt, err := db.Prepare(`
//--таблица статей
//	CREATE TABLE IF NOT EXISTS articles(
//		id INTEGER PRIMARY KEY,
//		title TEXT NOT NULL UNIQUE,
//		text TEXT NOT NULL);
//	CREATE INDEX IF NOT EXISTS idx_theme ON articles(title);
//	CREATE TABLE IF NOT EXISTS comments(
//		id INTEGER PRIMARY KEY,
//		text TEXT NOT NULL,
//		score INTEGER);
//`)
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	//обязательно закрываем, чтобы освободить ресурсы
//	defer stmt.Close()
//
//	_, err = stmt.Exec()
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	//Таблица комментариев-----------------------
//	stmt, err = db.Prepare(`
//	--таблица комментариев
//	CREATE TABLE IF NOT EXISTS comments(
//		id INTEGER PRIMARY KEY,
//		text TEXT NOT NULL,
//		score INTEGER);
//	`)
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	_, err = stmt.Exec()
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	////заполняем текстами
//	//articles := []struct {
//	//	title string
//	//	text  string
//	//}{
//	//	{"title 1", "article 1"},
//	//	{"title 2", "article 2"},
//	//	{"title 3", "article 3"},
//	//	{"title 4", "article 4"},
//	//	{"title 5", "article 5"},
//	//	{"title 6", "article 6"},
//	//	{"title 7", "article 7"},
//	//	{"title 8", "article 8"},
//	//}
//	//stmt, err = db.Prepare(`INSERT INTO articles (id, title, text) VALUES (?,?,?)`)
//	//if err != nil {
//	//	return nil, fmt.Errorf("%s: %w", op, err)
//	//}
//	//
//	//for id, article := range articles {
//	//	if _, err = stmt.Exec(id+1, article.title, article.text); err != nil {
//	//		return nil, fmt.Errorf("%s: %w", op, err)
//	//	}
//	//}
//	//
//	////заполняем комментариями
//	//comments := []struct {
//	//	text  string
//	//	score int64
//	//}{
//	//	{"comment 1", 50},
//	//	{"comment 2", 51},
//	//	{"comment 3", 52},
//	//	{"comment 4", 53},
//	//	{"comment 5", 54},
//	//	{"comment 6", 55},
//	//	{"comment 7", 56},
//	//	{"comment 8", 57},
//	//}
//	//stmt, err = db.Prepare(`INSERT INTO COMMENTS (id, text, score) VALUES (?,?,?)`)
//	//if err != nil {
//	//	return nil, fmt.Errorf("%s: %w", op, err)
//	//}
//	//
//	//for id, comment := range comments {
//	//	if _, err = stmt.Exec(id+1, comment.text, comment.score); err != nil {
//	//		return nil, fmt.Errorf("%s: %w", op, err)
//	//	}
//	//}
//
//	return &Storage{db: db}, nil
//}

// SaveURL сохранить
//func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
//	const op = "storage.sqlite.SaveURL"
//
//	// Подготавливаем запрос (проверка корректности синтаксиса)
//	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?, ?)")
//	if err != nil {
//		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
//	}
//
//	//выполняем запрос
//	res, err := stmt.Exec(urlToSave, alias)
//	if err != nil {
//		// Здесь мы приводим полученную ошибку ко внутреннему типу библиотеки sqlite3,
//		// чтобы посмотреть, не является ли эта ошибка sqlite3.ErrConstraintUnique.
//		// Если это так, значит, мы попытались добавить дубликат имеющейся записи. Об этом мы сообщим в вызывающую функцию, вернув уже свою ошибку для данной ситуации: storage.ErrURLExists. Получив ее, сервер сможет сообщить клиенту о том, что такой alias у нас уже есть.
//		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
//			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
//		}
//
//		var e sqlite3.Error
//		fmt.Println(e)
//
//		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
//	}
//
//	id, err := res.LastInsertId()
//	if err != nil {
//		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
//	}
//
//	//Возвращаем ID
//	return id, nil
//}

//func (s *Storage) makeStructJSON(queryText string, w http.ResponseWriter) error {
//
//	// returns rows *sql.Rows
//	rows, err := s.db.Query(queryText)
//	if err != nil {
//		return err
//	}
//	columns, err := rows.Columns()
//	if err != nil {
//		return err
//	}
//
//	count := len(columns)
//	values := make([]interface{}, count)
//	scanArgs := make([]interface{}, count)
//	for i := range values {
//		scanArgs[i] = &values[i]
//	}
//
//	masterData := make(map[string][]interface{})
//
//	for rows.Next() {
//		err := rows.Scan(scanArgs...)
//		if err != nil {
//			return err
//		}
//		for i, v := range values {
//
//			x := v.([]byte)
//
//			//NOTE: FROM THE GO BLOG: JSON and GO - 25 Jan 2011:
//			// The json package uses map[string]interface{} and []interface{} values to store arbitrary JSON objects and arrays; it will happily unmarshal any valid JSON blob into a plain interface{} value. The default concrete Go types are:
//			//
//			// bool for JSON booleans,
//			// float64 for JSON numbers,
//			// string for JSON strings, and
//			// nil for JSON null.
//
//			if nx, ok := strconv.ParseFloat(string(x), 64); ok == nil {
//				masterData[columns[i]] = append(masterData[columns[i]], nx)
//			} else if b, ok := strconv.ParseBool(string(x)); ok == nil {
//				masterData[columns[i]] = append(masterData[columns[i]], b)
//			} else if "string" == fmt.Sprintf("%T", string(x)) {
//				masterData[columns[i]] = append(masterData[columns[i]], string(x))
//			} else {
//				fmt.Printf("Failed on if for type %T of %v\n", x, x)
//			}
//
//		}
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//
//	err = json.NewEncoder(w).Encode(masterData)
//
//	if err != nil {
//		return err
//	}
//	return err
//}

func (s *Storage) GetArticleData(id string) ([]models.ArticleInfo, error) {
	const op = "storage.sqlite.GetArticleData"

	var result []models.ArticleInfo

	if err := s.db.Select(&result, "SELECT id, title, text, 0 as rating FROM articles"); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", result)

	return result, nil
}

//func (s *Storage) GetArticleData_Tst(id string) ([]models.ArticleInfo, error) {
//	const op = "storage.sqlite.GetArticleData"

//row := s.db.QueryRow(`
//SELECT id, title, text, 0 as rating FROM articles WHERE id = ?`, id)
//
//// Parse row into ArticleInfo struct
//var result []models.ArticleInfo
////activity := api.Activity{}
//var err error
//if err = row.Scan(&result.id activity.ID, &activity.Time, &activity.Description);
//	err == sql.ErrNoRows {
//	log.Printf("Id not found")
//	return api.Activity{}, ErrIDNotFound
//}

//// Подготавливаем запрос (проверка корректности синтаксиса)
//stmt, err := s.db.Prepare("SELECT id, title, text, 0 as rating FROM articles WHERE id = ?")
//if err != nil {
//	return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
//}
//
//var result []models.ArticleInfo
//
////в параметрах используем указатель, чтобы получить результаты
//err = stmt.QueryRow(id).Scan(&result)
////err = stmt.QueryRow("SELECT id, title, text, 0 as rating FROM articles WHERE id IN (2,3)").Scan(&result)
//fmt.Println("result:", result)
////если строки не найдено - возвращаем пустую строку
//if errors.Is(err, sql.ErrNoRows) {
//	return nil, storage.ErrDataNotFound
//} else {
//	fmt.Println("its ok. len=", len(result))
//}
//
//if err != nil {
//	return nil, fmt.Errorf("%s: execute statement: %w", op, err)
//}

//	return result, nil
//}

func (s *Storage) GetArticleData_OLD(id string) ([]models.ArticleInfo, error) {
	const op = "storage.sqlite.GetArticleData"

	// Подготавливаем запрос (проверка корректности синтаксиса)
	stmt, err := s.db.Prepare("SELECT id, title, text, 0 as rating FROM articles WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var result []models.ArticleInfo

	//в параметрах используем указатель, чтобы получить результаты
	err = stmt.QueryRow(id).Scan(&result)
	//err = stmt.QueryRow("SELECT id, title, text, 0 as rating FROM articles WHERE id IN (2,3)").Scan(&result)
	fmt.Println("result:", result)
	//если строки не найдено - возвращаем пустую строку
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrDataNotFound
	} else {
		fmt.Println("its ok. len=", len(result))
	}

	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
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

// DeleteURL Удалить запись из БД по алиасу
//func (s *Storage) DeleteURL(alias string) error {
//	const op = "storage.sqlite.DeleteURL"
//
//	// Подготавливаем запрос (проверка корректности синтаксиса)
//	stmt, err := s.db.Prepare("DELETE FROM url WHERe alias = ?")
//	if err != nil {
//		return fmt.Errorf("%s: prepare statement: %w", op, err)
//	}
//
//	//выполняем запрос
//	_, err = stmt.Exec(alias)
//	if err != nil {
//		return fmt.Errorf("%s: execute statement: %w", op, err)
//	}
//
//	return nil
//}
