package models

type ArticleInfo struct {
	Id     int64    `db:"id" json:"id"`
	Title  string   `db:"title" json:"title"`
	Text   string   `db:"text" json:"text"`
	Rating *float64 `db:"rating" json:"rating"` // Используем указатель, поскольку в БД может быть значение null
}
