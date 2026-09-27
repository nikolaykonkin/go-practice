# reflect-validation

Валидатор структур на Go через рефлексию. Проверяет поля по тегам `validate` (`min`, `max`, `regexp`) и возвращает первое несоответствие. Демонстрирует применение пакета `reflect` для универсального кода.

## Что показывает

- обход полей структуры через `reflect.Type` и `reflect.Value`;
- чтение тегов через `field.Tag.Get("validate")`;
- парсинг нескольких правил у одного поля (`min=3;max=10`);
- поддержку `min`/`max` для строк (по длине) и чисел (по значению);
- поддержку `regexp` для строк;
- корректный подсчёт длины строки в рунах (кириллица, иероглифы);
- обработку указателей и `nil`.

## Запуск

```bash
go run main.go
```

## Тесты

```bash
go test ./... -v
```

## Пример использования

```go
type User struct {
	Name  string `validate:"min=3"`
	Age   int    `validate:"min=18;max=65"`
	Email string `validate:"regexp=^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
}

if err := Validate(User{Name: "Ив", Age: 18, Email: "test@example.com"}); err != nil {
	fmt.Println("Validation error:", err)
}
// Validation error: поле Name: длина 2 меньше минимальной 3
```

## Пример вывода

```
Тест 1: ошибка: поле Name: длина 2 меньше минимальной 3
Тест 2: ошибка: поле Age: значение 70 больше максимального 65
Тест 3: ошибка: поле Email: значение "invalid email" не соответствует формату
Тест 4: ок
```

## Структура

```
reflect-validation/
├── main.go
├── main_test.go
├── go.mod
└── README.md
```