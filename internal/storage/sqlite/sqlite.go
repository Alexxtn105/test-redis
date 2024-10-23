// internal/storage/sqlite/sqlite.go

package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"test-redis/internal/storage"
)

// Storage Структура объекта Storage
type Storage struct {
	db *sql.DB //из пакета "database/sql"
}

// NewStorage Конструктор объекта Storage
func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок

	// Подключаемся к БД
	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: можно прикрутить миграции, для тренировки
	// создаем таблицу, если ее еще нет
	stmt, err := db.Prepare(`
--таблица статей
	CREATE TABLE IF NOT EXISTS articles(
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL UNIQUE,
		text TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_theme ON articles(title);
	CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY,
		text TEXT NOT NULL,
		score INTEGER);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//обязательно закрываем, чтобы освободить ресурсы
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//Таблица комментариев-----------------------
	stmt, err = db.Prepare(`
--таблица комментариев
	CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY,
		text TEXT NOT NULL,
		score INTEGER);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

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

// GetData - получить ссылку по ее алиасу
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
