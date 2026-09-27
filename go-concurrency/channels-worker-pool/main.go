package main

import (
	"fmt"
	"time"
)

// Order описывает заказ в интернет-магазине
type Order struct {
	ID       int
	Customer string
}

// generateOrders отправляет 10 заказов в канал и закрывает его
// Закрытие канала — сигнал воркерам, что новых заказов не будет
func generateOrders(out chan<- Order) {
	for i := 1; i <= 10; i++ {
		order := Order{
			ID:       i,
			Customer: fmt.Sprintf("Customer_%d", i),
		}
		fmt.Printf("Создан заказ %d\n", order.ID)
		out <- order
	}
	close(out)
}

// worker читает заказы из канала, обрабатывает их и отправляет сигнал о завершении в канал done
// Цикл for range завершается автоматически, когда канал orders закрыт и все заказы прочитаны
func worker(id int, orders <-chan Order, done chan<- bool) {
	for order := range orders {
		fmt.Printf("Worker %d начал обработку заказа %d\n", id, order.ID)

		// Имитация обработки заказа
		time.Sleep(500 * time.Millisecond)

		fmt.Printf("Worker %d завершил заказ %d (%s)\n", id, order.ID, order.Customer)
	}

	fmt.Printf("Worker %d завершил работу\n", id)
	done <- true
}

func main() {
	// Небуферизированный канал для заказов — по условию задания
	orders := make(chan Order)
	// Канал для сигналов завершения воркеров
	done := make(chan bool)

	const workerCount = 3

	// Запускаем воркеров
	for i := 1; i <= workerCount; i++ {
		go worker(i, orders, done)
	}

	// Генерация заказов в отдельной горутине
	go generateOrders(orders)

	// Ждем завершения всех воркеров
	for i := 0; i < workerCount; i++ {
		<-done
	}

	fmt.Println("Все заказы обработаны, программа завершена")
}
