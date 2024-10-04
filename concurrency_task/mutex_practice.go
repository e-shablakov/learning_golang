package main

import (
	"fmt"
	"sync"
)

var ArrayNumber []int = []int{1, 2, 3}

func sumArray(arr []int, sumC chan int, maxVal *int, mtx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for _, v := range arr {
		sum += v

		mtx.Lock()
		if v > *maxVal {
			*maxVal = v
		}
		mtx.Unlock()
	}
	sumC <- sum
}

func main() {
	sumCh := make(chan int, 2)
	var maxVal int
	var mtx sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)
	go sumArray(ArrayNumber[:len(ArrayNumber)/2], sumCh, &maxVal, &mtx, &wg)
	go sumArray(ArrayNumber[len(ArrayNumber)/2:], sumCh, &maxVal, &mtx, &wg)

	go func() {
		wg.Wait()
		close(sumCh)
	}()

	x, y := <-sumCh, <-sumCh
	fmt.Println("Sums:", x, y, "Total:", x+y)

	fmt.Println("Max value:", maxVal)
}
