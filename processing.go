package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func ProcessSequentially(numbers []int, applier func(int) int) []int {
	result := make([]int, len(numbers))
	for index, num := range numbers {
		result[index] = applier(num)
	}
	return result
}

func ProcessSimultaneously(numbers []int, applier func(int) int, threads int) []int {
	data := make([]int, len(numbers))
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	chunkSize := int(math.Ceil(float64(len(numbers)) / float64(threads)))
	for i := 0; i < threads; i++ {
		start := i * chunkSize
		end := start + chunkSize

		if start > len(numbers)-1 {
			break
		}

		if end > len(numbers) {
			end = len(numbers)
		}

		wg.Add(1)
		go func(start int, end int) {
			defer wg.Done()
			buffer := ProcessSequentially(data[start:end], applier)
			mutex.Lock()
			copy(data[start:end], buffer)
			mutex.Unlock()
		}(start, end)
	}
	wg.Wait()
	return data
}

func CountExecutionTime(action func(), text string) string {
	start := time.Now()
	action()
	time := time.Since(start)
	return fmt.Sprintf("%s execution took %v", text, time)
}
