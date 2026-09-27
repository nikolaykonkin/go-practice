# tasks-api-crud

REST API для управления списком задач (CRUD) на Go. Демонстрирует
создание HTTP-сервиса с JSON-ответами, потокобезопасным in-memory
хранилищем, логированием и слоистой архитектурой.

## Что показывает

- CRUD-эндпоинты на стандартной библиотеке `net/http` без сторонних роутеров;
- in-memory хранилище с защитой через `sync.RWMutex`;
- единый формат JSON-ответов и ошибок (`{"error": "..."}`);
- корректные HTTP-статусы (200/201/204/400/404/405/500);
- логирование запросов с методом, путём, статусом и временем обработки;
- слоистую архитектуру: `cmd` / `internal/handlers` / `internal/models` / `internal/storage`.

## Запуск

```bash
go run ./cmd/server
```

Сервер слушает `http://localhost:8080`. Проверка:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## API эндпоинты

| Метод  | Путь          | Описание                    | Успешный код | Ошибки          |
| ------ | ------------- | ---------------------------- | ------------ | --------------- |
| GET    | `/tasks`      | Получить список всех задач   | 200          | —                |
| POST   | `/tasks`      | Создать новую задачу         | 201          | 400              |
| GET    | `/tasks/{id}` | Получить задачу по ID        | 200          | 404              |
| PUT    | `/tasks/{id}` | Обновить задачу целиком      | 200          | 400, 404         |
| DELETE | `/tasks/{id}` | Удалить задачу                | 204          | 404              |
| GET    | `/health`     | Проверка состояния сервиса   | 200          | —                |

## Формат данных

```json
{
  "id": 1,
  "title": "Купить молоко",
  "done": false,
  "created_at": "2026-08-14T10:30:00Z"
}
```

- `id` — генерируется на сервере, значение из тела запроса игнорируется.
- `title` — обязательное поле, не может быть пустым или состоять только из пробелов.
- `done` — статус выполнения, по умолчанию `false`.
- `created_at` — генерируется на сервере в формате RFC3339, клиентское значение игнорируется.

## Формат ошибок

```json
{ "error": "сообщение об ошибке" }
```

| Код | Когда возвращается |
| --- | --- |
| 400 | Некорректный JSON, пустой `title`, невалидный `id` в пути, лишние данные после JSON |
| 404 | Задача с указанным `id` не найдена |
| 405 | Метод не поддерживается для данного пути (с заголовком `Allow`) |
| 500 | Внутренняя ошибка сервера |

## Примеры запросов

### Создать задачу

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Купить молоко"}'
```
```json
{"id":1,"title":"Купить молоко","done":false,"created_at":"2026-08-14T10:30:00Z"}
```

### Получить список задач

```bash
curl http://localhost:8080/tasks
```
```json
[{"id":1,"title":"Купить молоко","done":false,"created_at":"2026-08-14T10:30:00Z"}]
```

### Получить задачу по ID

```bash
curl http://localhost:8080/tasks/1
```
```json
{"id":1,"title":"Купить молоко","done":false,"created_at":"2026-08-14T10:30:00Z"}
```

### Обновить задачу

```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Купить молоко и хлеб", "done": true}'
```
```json
{"id":1,"title":"Купить молоко и хлеб","done":true,"created_at":"2026-08-14T10:30:00Z"}
```

### Удалить задачу

```bash
curl -i -X DELETE http://localhost:8080/tasks/1
# HTTP/1.1 204 No Content
```

## Ручное тестирование

Ниже — сценарии для каждого эндпоинта. Для каждого обязательного эндпоинта
приведён один корректный и один ошибочный запрос.

### `POST /tasks` — создание задачи

Корректный:

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Купить молоко"}'
# HTTP/1.1 201 Created
# {"id":1,"title":"Купить молоко","done":false,"created_at":"..."}
```

Ошибочный (пустой title):

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": ""}'
# HTTP/1.1 400 Bad Request
# {"error":"title is required"}
```

### `GET /tasks` — список задач

Корректный:

```bash
curl -i http://localhost:8080/tasks
# HTTP/1.1 200 OK
# [{"id":1,"title":"Купить молоко",...}]
```

Ошибочный (неподдерживаемый метод):

```bash
curl -i -X PATCH http://localhost:8080/tasks
# HTTP/1.1 405 Method Not Allowed
# Allow: GET, POST
# {"error":"method not allowed"}
```

### `GET /tasks/{id}` — задача по ID

Корректный:

```bash
curl -i http://localhost:8080/tasks/1
# HTTP/1.1 200 OK
# {"id":1,"title":"Купить молоко",...}
```

Ошибочный (несуществующий ID):

```bash
curl -i http://localhost:8080/tasks/999
# HTTP/1.1 404 Not Found
# {"error":"task not found"}
```

Ошибочный (лишний сегмент в пути):

```bash
curl -i http://localhost:8080/tasks/1/anything
# HTTP/1.1 400 Bad Request
# {"error":"invalid task id"}
```

### `PUT /tasks/{id}` — обновление задачи

Корректный:

```bash
curl -i -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Купить молоко и хлеб", "done": true}'
# HTTP/1.1 200 OK
# {"id":1,"title":"Купить молоко и хлеб","done":true,...}
```

Ошибочный (несуществующий ID):

```bash
curl -i -X PUT http://localhost:8080/tasks/999 \
  -H "Content-Type: application/json" \
  -d '{"title": "Не существует"}'
# HTTP/1.1 404 Not Found
# {"error":"task not found"}
```

### `DELETE /tasks/{id}` — удаление задачи

Корректный:

```bash
curl -i -X DELETE http://localhost:8080/tasks/1
# HTTP/1.1 204 No Content
```

Ошибочный (повторное удаление):

```bash
curl -i -X DELETE http://localhost:8080/tasks/1
# HTTP/1.1 404 Not Found
# {"error":"task not found"}
```

### Дополнительно: `/health`

```bash
curl -i http://localhost:8080/health
# HTTP/1.1 200 OK
# {"status":"ok"}
```

### Дополнительно: некорректный JSON

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{invalid json'
# HTTP/1.1 400 Bad Request
# {"error":"invalid JSON body"}
```

## Структура

```
tasks-api-crud/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   └── tasks.go
│   ├── middleware/
│   │   └── middleware.go
│   ├── models/
│   │   └── task.go
│   └── storage/
│       ├── storage.go
│       └── memory.go
├── go.mod
└── README.md
```