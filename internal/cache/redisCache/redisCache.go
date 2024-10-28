//internal/cache/redisCache/redisCache.go

package redisCache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"test-redis/internal/models"
	"time"
)

// Cache Структура объекта Cache
type Cache struct {
	client *redis.Client
}

// NewCache Конструктор объекта Cache
func NewCache(address string, password string, db int) (*Cache, error) {
	const op = "cache.redisCache.NewCache" // Имя текущей функции для логов и ошибок

	// Подключаемся к redis
	client := redis.NewClient(&redis.Options{
		Addr:     address,  // "localhost:6379",
		Password: password, // ""
		DB:       db,       // 0
	})

	return &Cache{client: client}, nil
}

var ctx = context.Background()

// SetKey Установка ключа
func (c *Cache) SetKey(key string, value any, expiration time.Duration) error {
	err := c.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetCachedArticle Получение данных о статье из кеша Redis
func (c *Cache) SetCachedArticle() ([]models.ArticleInfo, error) {

	//err := c.client.Set(ctx, "key", "value", 0).Err()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(c.client)

	return nil, nil
}

// GetCachedArticle Получение данных о статье из кеша Redis
func (c *Cache) GetCachedArticle(key string) ([]models.ArticleInfo, error) {

	//err := c.client.Set(ctx, "key", "value", 0).Err()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(c.client)

	return nil, nil
}
