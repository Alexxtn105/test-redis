# Установка и запуск Redis
## Установка Redis в Docker (докер должен быть запущен в Windows)
(Взято из видео https://www.youtube.com/watch?v=QpBaA6B1U90)
```bash
docker run -d -p 6379:6379 redis
```
## Просмотр запущенных контейнеров в docker
```bash
docker container ps
```
Будет выведен список запущенных контейнеров с их ИД

## Запуск командной строки redis-cli с **дефолтным** IP и номером порта (пример для ИД c2d01cdf6598)
```bash
docker exec -it c2d01cdf6598 redis-cli
```

## Запуск командной строки redis-cli с **конкретным** IP и номером порта (пример для ИД c2d01cdf6598)
```bash
docker exec -it c2d01cdf6598 redis-cli -h 127.0.0.1 -p 6379
```

## Проверка в redis-cli:
```bash
ping
```
В ответ должно прийти:
PONG

## Выход из redis-cli
```bash
exit
```

## Остановка работающего контейнера (пример для ИД c2d01cdf6598)
```bash
docker container stop c2d01cdf6598
```

# Команды redis-cli
## SET - Установка значения
Пример:
```bash
SET firstKey "Hello"
SET secondKey "World"
SET num 1
```
## GET - Получение значения
Пример:
```bash
GET firstKey
GET secondKey
```

## DEL - Удаление значения
Пример:
```bash
DEL firstKey
DEL secondKey
```

## KEYS - просмотр всех значений
```bash
KEYS *
```

## INCRBY - увеличение численного значения
```bash
INCRBY num 3  
```

## LPUSH и RPUSH - добавление в список слева и справа соответственно
```bash
# добавляем в список слева
LPUSH cars Toyota
# еще раз добавляем слева
LPUSH cars BMW
# теперь справа
RPUSH cars KIA
```

## LRANGE - посмотреть срез списка (индекс с .. по)
```bash
# вывод всех данных в списке
LRANGE cars 0 -1    
```

## LPOP и RPOP - получить значение из списка (справа или слева) с удалением его из списка
```bash
# получение левого элемента списка
LPOP cars 
```

## LLEN - получить длину списка
```bash
# получение длины списка
LLEN cars 
```

## LMOVE - перемещение элемента из одного списка в другой
```bash
# получение длины списка
LMOVE cars sold LEFT LEFT
```


# Работа с хеш-таблицами
## HSET - создание записи в Хэш-таблице
```bash
# Пример создания записи в хеш-таблице
HSET iPhone brand Apple model "iPhone X" memory 64 color Black
```

## HGET - получение ОДНОГО значения из записи Хэш-таблицы
```bash
# Пример получения названия модели из записи выше
HGET iPhone model
# или бренда
HGET iPhone brand
# или объем памяти
HGET iPhone memory
```

## HMGET - получение НЕСКОЛЬКИХ значений из записи Хэш-таблицы
```bash
# Пример получения нескольких значений из записи Хэш-таблицы
HMGET iPhone model brand unknown
```

## HGETALL - получение ВСЕХ значений из записи Хэш-таблицы
```bash
# Пример получения нескольких значений из записи Хэш-таблицы
HGETALL iPhone 
```