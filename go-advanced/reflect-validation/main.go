package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Validate проверяет структуру v на соответствие правилам,
// заданным через теги "validate" (min, max, regexp)
// Возвращает nil, если проверка пройдена, либо ошибку
// с описанием первого найденного несоответствия
func Validate(v any) error {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if !val.IsValid() || typ == nil {
		return fmt.Errorf("validate: передан nil")
	}

	// Разворачиваем указатели (в том числе многократные)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Errorf("validate: передан nil-указатель")
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("validate: ожидалась структура, получено %s", val.Kind())
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		for _, rule := range strings.Split(tag, ";") {
			if rule == "" {
				continue
			}

			parts := strings.SplitN(rule, "=", 2)
			ruleName := parts[0]
			var ruleValue string
			if len(parts) > 1 {
				ruleValue = parts[1]
			}

			if err := applyRule(field.Name, fieldValue, ruleName, ruleValue); err != nil {
				return err
			}
		}
	}

	return nil
}

func applyRule(fieldName string, fv reflect.Value, ruleName, ruleValue string) error {
	switch ruleName {
	case "min":
		return checkMin(fieldName, fv, ruleValue)
	case "max":
		return checkMax(fieldName, fv, ruleValue)
	case "regexp":
		return checkRegexp(fieldName, fv, ruleValue)
	default:
		return nil
	}
}

func checkMin(fieldName string, fv reflect.Value, ruleValue string) error {
	switch fv.Kind() {
	case reflect.String:
		minLen, err := strconv.Atoi(ruleValue)
		if err != nil {
			return fmt.Errorf("поле %s: некорректное правило min", fieldName)
		}
		runeLen := len([]rune(fv.String()))
		if runeLen < minLen {
			return fmt.Errorf("поле %s: длина %d меньше минимальной %d", fieldName, runeLen, minLen)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		minVal, err := strconv.ParseInt(ruleValue, 10, 64)
		if err != nil {
			return fmt.Errorf("поле %s: некорректное правило min", fieldName)
		}
		if fv.Int() < minVal {
			return fmt.Errorf("поле %s: значение %d меньше минимального %d", fieldName, fv.Int(), minVal)
		}
	}
	return nil
}

func checkMax(fieldName string, fv reflect.Value, ruleValue string) error {
	switch fv.Kind() {
	case reflect.String:
		maxLen, err := strconv.Atoi(ruleValue)
		if err != nil {
			return fmt.Errorf("поле %s: некорректное правило max", fieldName)
		}
		runeLen := len([]rune(fv.String()))
		if runeLen > maxLen {
			return fmt.Errorf("поле %s: длина %d больше максимальной %d", fieldName, runeLen, maxLen)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		maxVal, err := strconv.ParseInt(ruleValue, 10, 64)
		if err != nil {
			return fmt.Errorf("поле %s: некорректное правило max", fieldName)
		}
		if fv.Int() > maxVal {
			return fmt.Errorf("поле %s: значение %d больше максимального %d", fieldName, fv.Int(), maxVal)
		}
	}
	return nil
}

func checkRegexp(fieldName string, fv reflect.Value, pattern string) error {
	if fv.Kind() != reflect.String {
		return nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("поле %s: некорректное регулярное выражение", fieldName)
	}
	if !re.MatchString(fv.String()) {
		return fmt.Errorf("поле %s: значение %q не соответствует формату", fieldName, fv.String())
	}
	return nil
}

type User struct {
	Name  string `validate:"min=3"`
	Age   int    `validate:"min=18;max=65"`
	Email string `validate:"regexp=^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
}

func main() {
	tests := []User{
		{Name: "Ив", Age: 18, Email: "test@example.com"},
		{Name: "Иван", Age: 70, Email: "test@example.com"},
		{Name: "Иван", Age: 35, Email: "invalid email"},
		{Name: "Иван", Age: 35, Email: "test@example.com"},
	}

	for i, u := range tests {
		if err := Validate(u); err != nil {
			fmt.Printf("Тест %d: ошибка: %v\n", i+1, err)
		} else {
			fmt.Printf("Тест %d: ок\n", i+1)
		}
	}
}
