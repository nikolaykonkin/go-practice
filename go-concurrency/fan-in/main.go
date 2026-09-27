package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Metric описывает одну метрику от источника данных
type Metric struct {
	Source string    // "CPU", "Memory" или "Network"
	Value  float64   // значение метрики
	Time   time.Time // время создания метрики
}

// cpuMetrics эмулирует метрики загрузки процессора:
// 5 значений с интервалом 800 мс, диапазон 0–100
func cpuMetrics() <-chan Metric {
	ch := make(chan Metric)
	go func() {
		defer close(ch)
		for i := 0; i < 5; i++ {
			ch <- Metric{
				Source: "CPU",
				Value:  rand.Float64() * 100,
				Time:   time.Now(),
			}
			time.Sleep(800 * time.Millisecond)
		}
	}()
	return ch
}

// memoryMetrics эмулирует метрики использования памяти:
// 5 значений с интервалом 1200 мс, диапазон 0–16384 МБ
func memoryMetrics() <-chan Metric {
	ch := make(chan Metric)
	go func() {
		defer close(ch)
		for i := 0; i < 5; i++ {
			ch <- Metric{
				Source: "Memory",
				Value:  rand.Float64() * 16384,
				Time:   time.Now(),
			}
			time.Sleep(1200 * time.Millisecond)
		}
	}()
	return ch
}

// networkMetrics эмулирует метрики сетевой активности:
// 5 значений с интервалом 1500 мс, диапазон 0–1000 Мбит/с
func networkMetrics() <-chan Metric {
	ch := make(chan Metric)
	go func() {
		defer close(ch)
		for i := 0; i < 5; i++ {
			ch <- Metric{
				Source: "Network",
				Value:  rand.Float64() * 1000,
				Time:   time.Now(),
			}
			time.Sleep(1500 * time.Millisecond)
		}
	}()
	return ch
}

// fanIn объединяет несколько каналов в один
// Для каждого источника запускается отдельная горутина:
// она читает метрики из своего канала и пишет их в общий
// Общий канал закрывается после завершения всех источников
func fanIn(channels ...<-chan Metric) <-chan Metric {
	out := make(chan Metric)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan Metric) {
			defer wg.Done()
			for m := range c {
				out <- m
			}
		}(ch)
	}

	// Отдельная горутина ждет завершения всех источников и закрывает выходной канал —
	// это позволяет main завершить цикл for range корректно
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	fmt.Println("Запуск системы мониторинга")

	cpuCh := cpuMetrics()
	memCh := memoryMetrics()
	netCh := networkMetrics()

	combined := fanIn(cpuCh, memCh, netCh)

	for metric := range combined {
		fmt.Printf("Источник: %s, Значение: %.2f, Время: %s\n",
			metric.Source, metric.Value, metric.Time.Format("15:04:05.000"))
	}

	fmt.Println("Мониторинг завершён")
}
