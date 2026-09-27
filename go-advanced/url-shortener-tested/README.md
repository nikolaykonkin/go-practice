# url-shortener-tested

HTTP-сервис сокращения URL на Go с полным покрытием тестами. Демонстрирует table-driven тесты, `httptest` и разделение бизнес-логики и HTTP-слоя.

## Что показывает

- HTTP-сервис с двумя эндпоинтами (`POST /shorten`, `GET /{short_id}`);
- in-memory хранилище с `sync.RWMutex`;
- валидацию URL через `url.ParseRequestURI`;
- генерацию ID через `crypto/rand` + URL-safe Base64;
- обработку коллизий и проброс ошибок генератора;
- table-driven тесты для бизнес-логики и `httptest` для обработчиков;
- интеграционный тест маршрутизации через `http.ServeMux`;
- тест на конкурентный доступ (запускать с `-race`).

## Запуск

```bash
go run .
```

Сервер слушает `http://localhost:8080`.

## Тесты

```bash
go test ./... -v       # все тесты с подробным выводом
go test ./... -race    # проверка на гонки данных
go test -cover ./...   # отчет о покрытии
```

Table-driven тесты бизнес-логики (`shortener_test.go`) и HTTP-обработчиков через `httptest` (`handlers_test.go`). Покрытие бизнес-логики — выше 80%. Есть интеграционный тест маршрутизации и тест на конкурентный доступ.

## API

### `POST /shorten`

Запрос:

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very/long/path"}'
```

Ответ:

```json
{"short_url":"abc123","original_url":"https://example.com/very/long/path"}
```

### `GET /{short_id}`

```bash
curl -i http://localhost:8080/abc123
# HTTP/1.1 302 Found
# Location: https://example.com/very/long/path
```

Если ID не найден — `404 Not Found` с телом `{"error":"url не найден"}`.

## Ручное тестирование

Запуск сервера:

```bash
go run .
# Сервер запущен на :8080
```

Создание короткой ссылки:

```bash
curl -i -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very/long/path"}'
# HTTP/1.1 200 OK
# Content-Type: application/json
# {"short_url":"...","original_url":"https://example.com/very/long/path"}
```

Редирект:

```bash
curl -i http://localhost:8080/<short_id_из_предыдущего_ответа>
# HTTP/1.1 302 Found
# Location: https://example.com/very/long/path
```

Ошибки:

```bash
# Некорректный JSON
curl -i -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{invalid'
# HTTP/1.1 400 Bad Request
# {"error":"некорректный JSON"}

# Невалидный URL
curl -i -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "not-a-url"}'
# HTTP/1.1 400 Bad Request
# {"error":"невалидный URL"}

# Несуществующий short_id
curl -i http://localhost:8080/doesnotexist
# HTTP/1.1 404 Not Found
# {"error":"url не найден"}

# Неподдерживаемый метод
curl -i -X DELETE http://localhost:8080/shorten
# HTTP/1.1 405 Method Not Allowed
# Allow: POST
```

## Структура

```
url-shortener-tested/
├── main.go                — HTTP-сервер и обработчики
├── shortener.go           — бизнес-логика (URLShortener, генерация ID)
├── shortener_test.go      — unit-тесты бизнес-логики
├── handlers_test.go       — тесты HTTP-обработчиков (httptest)
├── go.mod
└── README.md
```