//internal/cache/redisCache/redisCache.go

package redisCache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"test-redis/internal/models"
	"time"
)

// ctx Текущий контекст
var ctx = context.Background()

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

// SetKey Установка ключа
func (c *Cache) SetKey(id string, value any, expiration time.Duration) error {
	err := c.client.Set(ctx, "article:"+id, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// SetCachedArticle Получение данных о статье из кеша Redis
func (c *Cache) SetCachedArticle(id string, value any) error {

	err := c.client.Set(ctx, "article:"+id, value, time.Second*100).Err()

	if err != nil {
		fmt.Println("error setting value ", err)
		return err
	}

	return nil
}

// GetCachedArticle Получение данных о статье из кеша Redis
func (c *Cache) GetCachedArticle(id string) ([]models.ArticleInfo, error) {
	raw, err := c.client.Get(ctx, "article:"+id).Result()
	fmt.Println(raw)
	if err == redis.Nil {
		return nil, fmt.Errorf("key %s does not exist", "article:"+id)
	} else if err != nil {
		return nil, err
	}
	var info []models.ArticleInfo
	err = json.Unmarshal([]byte(raw), &info)
	fmt.Println("got from cache:")
	fmt.Println(info)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal raw value %s", "article:"+id)
	}

	return info, nil
}

func (c *Cache) GetCachedArticleAsString(key string) (string, error) {
	raw, err := c.client.Get(ctx, "article:"+key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist", "article:"+key)
	} else if err != nil {
		return "", err
	}
	return raw, nil
}
