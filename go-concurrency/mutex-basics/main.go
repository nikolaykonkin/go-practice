package main

import (
	"fmt"
	"sync"
)

func main() {
	// Общая переменная-счетчик, которую инкрементируют все горутины
	var counter int

	// Mutex делает инкремент counter атомарным: пока одна горутина
	// держит блокировку, другие ждут освобождения
	var mu sync.Mutex

	// WaitGroup отслеживает завершение всех горутин
	var wg sync.WaitGroup

	const goroutines = 10
	const increments = 1000

	// wg.Add(n) сообщает WaitGroup, что нужно дождаться n горутин
	// Вызывается до запуска горутин — иначе wg.Wait может вернуться раньше,
	// чем счетчик достигнет n
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < increments; j++ {
				// Блокируем counter, чтобы предотвратить одновременное
				// изменение несколькими горутинами
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	// wg.Wait блокирует main до тех пор, пока счетчик WaitGroup
	// не станет нулем — то есть пока все горутины не вызовут wg.Done
	wg.Wait()

	fmt.Println("Итоговое значение счётчика:", counter)
	fmt.Println("Ожидаемое значение:", goroutines*increments)
}
