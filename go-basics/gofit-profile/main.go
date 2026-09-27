package main

import "fmt"

// Название программы
const AppName = "GoFit"

func main() {
	// Профиль пользователя
	var (
		name        string  = "Иван"
		age         int     = 30
		height      float64 = 1.75
		isSubscribed bool   = true
	)

	fmt.Printf("Добро пожаловать в %s!\n", AppName)
	fmt.Println("Профиль пользователя:")
	fmt.Printf("Имя: %s\n", name)
	fmt.Printf("Возраст: %d лет\n", age)
	fmt.Printf("Рост: %.2f м\n", height)
	fmt.Printf("Подписан на рассылку: %t\n", isSubscribed)
}
