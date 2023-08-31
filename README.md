# AvitoUserSegmenting
Cервис, который хранит пользователя и его сегменты, в которых он состоит.
Дает возможность создания и удаления сегментов, а также добавление и удаление пользователей в сегмент.  
Используемые технологии:  
- Echo (веб фреймворк)
- PosgreSQL (база данных)
- pgx (драйвер для работы с БД)
- Docker + Docker Compose (контейниризация микросервиса)
- migrate (миграции для БД)
- swag (генерация документации)
- logrus (логгер)

## Начнем!
Для подготовки сервиса нужно:
- Склонировать репозиторий:
```
https://github.com/korsilyn/AvitoUserSegmenting
```
- Изменить environment переменные в docker-compose.yaml под себя
- Настроить config/config.yaml

## Запуск
```
make compose-up
```

## Примеры

Некоторые примеры запросов
- [Создание сегмента](#create-seg)
- [Удаление сегмента](#del-seg)
- [Добавление/Удаление сегментов](#add-remove)
- [Получение списка сегментов](#seg-list)

### Создание сегмента <a name="create-seg"></a>

При создании сегмента реализована опция указания процента пользователей (из общего колличества), которые попадут в этот сегмент автоматически, а так же есть возможность установить TTL:
```curl
curl --location --request POST 'http://localhost:8080/segment' \
--header 'Content-Type: application/json' \
--data-raw '{
    "slug": "AVITO_SALE_60",
    "expiration_date": "2023-12-31T23:59:59Z",
    "random_percentage": 0.0
}'
```
Пример ответа:
```json
{
   "message": "Segment and user assignments created successfully"
}
```

### Удаление сегмента <a name="del-seg"></a>

Удаление сегмента по указанному slug:
```curl
curl --location --request DELETE 'http://localhost:8080/segment' \
--header 'Content-Type: application/json' \
--data-raw '{
    "slug": "AVITO_SALE_60"
}'
```
Пример ответа:
```json
{
   "message": "Segment deleted successfully",
   "segment_id": 1
}
```

### Добавление/Удаление сегментов <a name="add-remove"></a>

Добавление / удаление сегментов пользователя списком без перетирания существующих сегментов с возможностью установить TTL.
```curl
curl --location --request POST 'http://localhost:8080/user/segments' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id": 1,
    "add": [
    {
      "slug": "AVITO_SALE_10",
      "expiration_date": "2023-12-31T23:59:59Z"
    },
    {
      "slug": "AVITO_SALE_30",
      "expiration_date": "2023-11-30T23:59:59Z"
    },
    {
      "slug": "AVITO_SALE_20",
      "expiration_date": "2023-11-30T23:59:59Z"
    }
    ], 
    "remove": ["AVITO_SALE_40"]
}'
```
Пример ответа:
```json
{
   "message": "User segments updated successfully",
   "user_id": 1
}

```

### Получение списка сегментов <a name="seg-list"></a>

Получение списка сегментов пользователя по id:
```curl
curl --location --request GET 'http://localhost:8080/user/segments' \
--header 'Content-Type: application/json' \
--data-raw '{
   "user_id": 1
}'
```
Пример ответа:
```json
{
   "segments": ["AVITO_SALE_10","AVITO_SALE_30"],
   "user_id": 1
}
```

## Недоработки (почти все мог доделать, было бы чуть больше времени, планирую доделать часть после дедлайна просто для себя)
- Нету тестов вообще
- Не сделана выдача в CSV (1 доп), хоть и были заготовки
- Возможно где-то не очень хороший код, где то можно было сделать получше, поаккуратнее, где то не хватает проверок данных и т.п.
Простой пример - дубликаты убираются не при вносе пользователя, а при выдаче, могло уменьшить запросы к БД.
- Сегменты хранятся по отдельности, в связке пользователь - сегмент 1 ряд. По хорошему так не делать, потому что при тысячах юзеров
запросы будут достаточно долгие. Идея как исправить появилась достаточно поздно.
- При запросе времени из БД выдается UTC, а отправляется MSK+3

## Вопросы
1. Встал вопрос о хранении пользователей без сегментов и в принципе создания БД под юзеров. Подумал и 
решил, что стоит не проверять ID у юзеров, по факту это должна быть информация с внешней БД.
Если у пользователя нет сегментов - просто не храним, иначе храним с временем удаления из сегмента
и при запросе смотрим, чтобы время было актуальным. Заодно удобно будет с выдачей за месяц.  
2. Доп задание 3. "10% пользователей попадают в сегмент". Нет указания на то, откуда брать
id пользователей и общее количество пользователей. Для данного задания решил сделать возможные 
id от 1000 до 1099 (100 юзеров). По хорошему, на проде забирать эти данные с основной БД или 
придумать что нибудь похожее.  
3. В задании написано, что добавление и удаление сегментов должно обрабатываться одним методом, но я посчитал, 
что будет удобнее читать и исправлять код, имея в API методы /add и /remove, чем 1 большой метод /change.  
4. Если пользователю второй раз добавили сегмент, нужно ли обновлять время добавления? 
Я решил не обновлять, потому что такие ситуации мне кажется будут очень редкими.  
5. Доп задание 2. В чем получать TTL и как его хранить? Я решил сразу объединить 2 и 3 
допы, и сразу заносить TTL в поле removed_at. При получении списка пользователей появилась 
доп проверка. TTL передается в часах, просто int цифрой, и потом обрабатывается через 
time.Duration.  
6. Нужна ли авторизация? В задании не указано, но мне кажется, для такого API  нужна либо 
авторизация, либо доступ только из локальной сетки. Тк у меня не осталось времени, 
я оставил этот вопрос просто на грани рассуждений.
