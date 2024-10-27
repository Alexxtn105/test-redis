# Тестовый проект с использованием redis

## ЗАПУСК СЕРВИСА:
```bash
go run ./cmd/test-redis/main.go --config=./config/local.yaml
```

Для запуска (в bash, с использованием переменной окружения):
```bash
CONFIG_PATH="./config/local.yaml" go run  "./cmd/test-redis/main.go"
```

ЗАПУСК ТЕСТОВ:
```bash
go test ./tests -count=1 -v
```

## ПРИМЕР РУЧНОЙ УСТАНОВКИ ТЕГА
```bash
git tag v0.0.1 && git push origin v0.0.1
```


Страница для тестирования API
https://jsonplaceholder.typicode.com/users