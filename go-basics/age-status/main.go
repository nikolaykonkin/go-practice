package main

import "fmt"

func main() {
	var age int

	for i := 1; i <= 5; i++ {
		fmt.Print("Введите возраст: ")
		fmt.Scan(&age)

		switch {
		case age < 0:
			fmt.Println("Ошибка: возраст не может быть отрицательным")
		case age < 18:
			fmt.Println("Ребёнок")
		case age <= 65:
			fmt.Println("Взрослый")
		default:
			fmt.Println("Пенсионер")
		}
	}
}
