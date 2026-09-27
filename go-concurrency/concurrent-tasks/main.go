package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task описывает задачу для выполнения
type Task struct {
	ID   int
	Name string
	Type string
}

// Result описывает результат выполнения задачи
type Result struct {
	TaskID   int
	TaskName string
	Success  bool
	Duration time.Duration
	Message  string
}

// simulateIOWork имитирует I/O-операцию со случайной задержкой
func simulateIOWork(task Task, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done()

	fmt.Printf("Начало I/O задачи %d: %s\n", task.ID, task.Name)
	start := time.Now()

	// Имитация задержки 1–3 секунды
	sleepTime := time.Duration(rand.Intn(3)+1) * time.Second
	time.Sleep(sleepTime)

	results <- Result{
		TaskID:   task.ID,
		TaskName: task.Name,
		Success:  true,
		Duration: time.Since(start),
		Message:  "I/O операция завершена",
	}
	fmt.Printf("Завершена I/O задача %d: %s\n", task.ID, task.Name)
}

// simulateComputeWork имитирует вычислительную задачу
func simulateComputeWork(task Task, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done()

	fmt.Printf("Начало вычислительной задачи %d: %s\n", task.ID, task.Name)
	start := time.Now()

	// Небольшая нагрузка: число Фибоначчи от 25 до 29
	n := rand.Intn(5) + 25
	_ = fibonacci(n)

	results <- Result{
		TaskID:   task.ID,
		TaskName: task.Name,
		Success:  true,
		Duration: time.Since(start),
		Message:  "Вычисление завершено",
	}
	fmt.Printf("Завершена вычислительная задача %d: %s\n", task.ID, task.Name)
}

// fibonacci вычисляет n-е число Фибоначчи рекурсивно
// Используется как искусственная вычислительная нагрузка
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// monitorProgress следит за результатами и раз в 2 секунды печатает прогресс
// Завершается, когда канал results закрыт
func monitorProgress(totalTasks int, results <-chan Result, done chan<- bool) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	completed := 0
	for {
		select {
		case res, ok := <-results:
			if !ok {
				done <- true
				return
			}
			completed++
			fmt.Printf("Результат: %s (%v)\n", res.TaskName, res.Duration)
		case <-ticker.C:
			progress := float64(completed) / float64(totalTasks) * 100
			fmt.Printf("Прогресс: %.0f%% (%d/%d)\n", progress, completed, totalTasks)
		}
	}
}

func main() {
	fmt.Println("=== Демонстрация горутин в Go ===")

	ioTasks := []Task{
		{1, "Загрузка данных", "IO"},
		{2, "Чтение файла", "IO"},
		{3, "Запрос к API", "IO"},
		{4, "Сохранение данных", "IO"},
	}

	computeTasks := []Task{
		{5, "Фибоначчи анализ", "COMPUTE"},
		{6, "Обработка данных", "COMPUTE"},
		{7, "Математическая модель", "COMPUTE"},
		{8, "Оптимизация алгоритма", "COMPUTE"},
	}

	allTasks := append(ioTasks, computeTasks...)
	totalTasks := len(allTasks)

	results := make(chan Result, totalTasks)
	monitorDone := make(chan bool)
	var wg sync.WaitGroup

	fmt.Printf("\nЗапускаю %d задач:\n", totalTasks)
	for _, task := range allTasks {
		fmt.Printf("  %s (ID: %d, Тип: %s)\n", task.Name, task.ID, task.Type)
	}

	go monitorProgress(totalTasks, results, monitorDone)

	fmt.Println("\nЗапуск горутин...")
	startTime := time.Now()

	for _, task := range ioTasks {
		wg.Add(1)
		go simulateIOWork(task, &wg, results)
	}
	for _, task := range computeTasks {
		wg.Add(1)
		go simulateComputeWork(task, &wg, results)
	}

	fmt.Printf("Запущено %d горутин\n", totalTasks)

	wg.Wait()
	close(results)
	<-monitorDone

	totalExecutionTime := time.Since(startTime)
	fmt.Println("\n=== Программа завершена ===")
	fmt.Printf("Общее время выполнения: %v\n", totalExecutionTime)
	fmt.Println("Все горутины успешно завершены")
}
