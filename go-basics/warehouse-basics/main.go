package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Глобальные константы для категорий товаров
const (
	CategoryElectronics = "Электроника"
	CategoryFood        = "Продукты"
	CategoryClothes     = "Одежда"
	MaxItems            = 100 // Максимальное количество товаров на складе
)

// Глобальная переменная для подсчета товаров
var totalItems int

func main() {
	// Разные способы объявления и инициализации переменных

	// Способ 1: объявление с явным типом и последующим присваиванием
	var itemName string
	itemName = "Смартфон"

	// Способ 2: объявление с явным типом и инициализацией
	var itemPrice float64 = 999.99

	// Способ 3: несколько переменных одного типа
	var minQuantity, maxQuantity int = 5, 20

	// Способ 4: тип выводится компилятором
	var isAvailable = true

	// Способ 5: короткое объявление (только внутри функций)
	quantity := 15

	// Способ 6: несколько переменных разных типов одновременно
	var (
		itemID     int64  = addNewItem(itemName, quantity)
		itemColor  string = "Черный"
		itemWeight float32
		dateAdded  time.Time = time.Now()
	)

	itemWeight = 0.3

	// Приведение типов: без float64 было бы целочисленное деление
	percentInStock := float64(quantity) / float64(MaxItems) * 100

	category := CategoryElectronics

	discountedPrice := calculateDiscount(itemPrice, 15)

	displayItemInfo(itemID, itemName, quantity, discountedPrice, isAvailable, category)

	fmt.Println("\nДополнительная информация:")
	fmt.Printf("Цвет товара: %s\n", itemColor)
	fmt.Printf("Вес товара: %.2f кг\n", itemWeight)
	fmt.Printf("Итоговая стоимость за единицу со скидкой: %.2f руб.\n", discountedPrice)
	fmt.Printf("Минимальное количество: %d, Максимальное количество: %d\n", minQuantity, maxQuantity)
	fmt.Printf("Процент от максимального количества на складе: %.1f%%\n", percentInStock)
	fmt.Printf("Дата добавления: %s\n", dateAdded.Format("02-01-2006"))

	// Демонстрация разных числовых типов
	var (
		shortValue   int8      = 127
		intValue     int       = 1000000
		uintValue    uint      = 10000
		floatValue   float32   = 123.456
		complexValue complex64 = 1 + 2i
	)

	fmt.Println("\nРазные числовые типы:")
	fmt.Printf("int8: %d\n", shortValue)
	fmt.Printf("int: %d\n", intValue)
	fmt.Printf("uint: %d\n", uintValue)
	fmt.Printf("float32: %f\n", floatValue)
	fmt.Printf("complex64: %v\n", complexValue)

	// Работа со строками и байтами (руны vs байты)
	productCode := "ЯЩИК-12345"
	fmt.Printf("\nКод товара: %s, длина: %d символов, длина: %d байт\n",
		productCode, len([]rune(productCode)), len(productCode))

	updateTotalItems(quantity)
	fmt.Printf("\nОбщее количество товаров на складе: %d\n", totalItems)
}

// addNewItem добавляет новый товар и возвращает его ID
// ID генерируется на основе Unix-времени и случайного числа
// В Go 1.20+ rand.Seed не требуется — math/rand инициализируется автоматически
func addNewItem(name string, qty int) int64 {
	newID := time.Now().Unix() + int64(rand.Intn(1000))

	fmt.Printf("--- Добавление товара: %s (кол-во: %d) ---\n", name, qty)

	return newID
}

// calculateDiscount возвращает цену со скидкой percent процентов
// Если процент вне диапазона [0, 100], возвращает исходную цену
func calculateDiscount(price float64, percent float64) float64 {
	if percent < 0 || percent > 100 {
		fmt.Println("Ошибка: некорректный процент скидки")
		return price
	}

	discountAmount := price * (percent / 100)
	finalPrice := price - discountAmount

	fmt.Printf("Применена скидка %.0f%%. Экономия: %.2f руб.\n", percent, discountAmount)
	return finalPrice
}

// displayItemInfo выводит информацию о товаре в читаемом виде
func displayItemInfo(id int64, name string, qty int, price float64, isAvailable bool, category string) {
	fmt.Println("=== Информация о товаре ===")
	fmt.Printf("ID: %d\n", id)
	fmt.Printf("Название: %s\n", name)
	fmt.Printf("Категория: %s\n", category)
	fmt.Printf("Количество: %d\n", qty)
	fmt.Printf("Цена: %.2f руб.\n", price)

	if isAvailable {
		fmt.Println("Статус: В наличии")
	} else {
		fmt.Println("Статус: Нет в наличии")
	}
}

// updateTotalItems увеличивает глобальный счетчик товаров и не дает ему превысить MaxItems
func updateTotalItems(qty int) {
	totalItems += qty

	if totalItems > MaxItems {
		fmt.Println("Предупреждение: превышено максимальное количество товаров на складе")
		totalItems = MaxItems
	}
}
