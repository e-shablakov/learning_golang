package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	ConsumerCount = 2
	ProducerCount = 2
)

var (
	stopSignal int32
	buffer     []Task
	mutex      sync.Mutex
	wg         sync.WaitGroup
)

type Task struct {
	id   int
	name string
}

func producer(id int) {
	// выполнится после выхода из цикла
	// сообщает о завершении работы горутины
	defer wg.Done()

	for {
		// проверка стоп сигнала
		if atomic.LoadInt32(&stopSignal) == 1 {
			fmt.Printf("Producer %d stopping...\n", id)
			return
		}

		task := Task{id: rand.Intn(100), name: fmt.Sprintf("Task from %d producer", id)}

		// критическая секция для доступа к буфферу
		mutex.Lock()
		buffer = append(buffer, task)
		fmt.Printf("Producer %d added: %v\n", id, task)
		mutex.Unlock()

		time.Sleep(time.Second * 2)
	}
}

func consumer(id int) {
	// выполнится после выхода из цикла
	// сообщает о завершении работы горутины
	defer wg.Done()
	for {
		// проверка стоп сигнала
		if atomic.LoadInt32(&stopSignal) == 1 {
			fmt.Printf("Consumer %d stopping...\n", id)
			return
		}

		// критическая секция для доступа к буфферу
		mutex.Lock()
		if len(buffer) > 0 {
			item := buffer[0]
			buffer = buffer[1:]
			fmt.Printf("Consumer %d consumed: %v\n", id, item)
		}
		mutex.Unlock()

		time.Sleep(time.Second * 2)
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= ProducerCount; i++ {
		wg.Add(1)
		go producer(i)
	}

	for i := 1; i <= ConsumerCount; i++ {
		wg.Add(1)
		go consumer(i)
	}

	fmt.Println("Press any key to stop...")
	_, _ = fmt.Scanln()

	atomic.StoreInt32(&stopSignal, 1)

	// ожидает выполнения всех горутин
	wg.Wait()

	fmt.Println("All producers and consumers stopped.")
}
