# go-practice

Учебные проекты по курсу «Go-разработчик с нуля» (Нетология), сгруппированные по трем блокам программы. Каждая подпапка — независимый проект со своим `go.mod` и README.

## Структура

### go-basics — Основы программирования на Go

- [employees-list](./go-basics/employees-list/) — структуры, срезы, интерфейсы
- [warehouse-basics](./go-basics/warehouse-basics/) — основные конструкции языка
- [age-status](./go-basics/age-status/) — условия и циклы
- [gofit-profile](./go-basics/gofit-profile/) — базовые типы данных
- [read-process-write](./go-basics/read-process-write/) — функции и работа с файлами
- [modules](./go-basics/modules/) — модули и пакеты

### go-concurrency — Многопоточность Go

- [mutex-basics](./go-concurrency/mutex-basics/) — sync.Mutex, sync.WaitGroup
- [channels-worker-pool](./go-concurrency/channels-worker-pool/) — worker pool на каналах
- [fan-in](./go-concurrency/fan-in/) — паттерн Fan-in
- [concurrent-tasks](./go-concurrency/concurrent-tasks/) — I/O vs compute, ticker

### go-advanced — Продвинутое изучение Go

- [reflect-validation](./go-advanced/reflect-validation/) — валидатор структур через reflect
- [tasks-api-crud](./go-advanced/tasks-api-crud/) — REST API, in-memory storage
- [url-shortener-tested](./go-advanced/url-shortener-tested/) — REST API + table-driven tests

## Запуск

Каждая подпапка — независимый проект. Инструкции по запуску — в README соответствующей подпапки:

```bash
cd go-basics/employees-list
go run main.go
```

Для проектов с тестами:

```bash
cd go-advanced/url-shortener-tested
go test ./... -v
```

## Требования

Go 1.26 или выше.