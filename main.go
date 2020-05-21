package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var wgWorkers sync.WaitGroup

// Expected result: when using workers, the sum of the numbers is done faster.
func main() {
	println("Number of processors: " + strconv.Itoa(runtime.NumCPU()))

	// Play with this!!
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create a slice with a million numbers
	var bigSlice []int
	for i := 1; i <= 1000000; i++ {
		bigSlice = append(bigSlice, i)
	}

	// Create timers
	var start time.Time
	var finish time.Time

	// Get elapsed time of the sum with workers
	start = time.Now()
	result := sumWithWorkers(bigSlice)
	finish = time.Now()
	fmt.Println("Time elapsed with workers: ", finish.Sub(start))
	fmt.Println("Result: ", result)

	// Get elapsed time of the sum without using workers
	start = time.Now()
	result = sumNumbersFromSlice(bigSlice)
	finish = time.Now()
	fmt.Println("Time elapsed without using workers: ", finish.Sub(start))
	fmt.Println("Result: ", result)
}

func sumWithWorkers(numberSlice []int) int64 {
	// Check the amount of CPUs to tell how many goroutines can be run in parallel
	goroutineAmount := runtime.NumCPU()
	wgWorkers.Add(goroutineAmount)

	// Get the amount of numbers to sum in each goroutine
	amountOfNumbersToSumInEachGoroutine := len(numberSlice) / goroutineAmount

	var sumTotal int64 = 0
	for i := 0; i < goroutineAmount; i++ {
		// Decide the section of the slice that this worker will work with
		start := amountOfNumbersToSumInEachGoroutine * i
		end := start + amountOfNumbersToSumInEachGoroutine
		sliceSectionToSum := numberSlice[start:end]

		// Create a goroutine that will sum all the numbers from slice and add the result to the total
		go func(sliceSectionToSum []int, sumTotal *int64) {
			defer wgWorkers.Done()
			atomic.AddInt64(sumTotal, sumNumbersFromSlice(sliceSectionToSum))
		}(sliceSectionToSum, &sumTotal)
	}

	// Wait until all the workers finish summing their sections to return the result
	wgWorkers.Wait()

	return sumTotal
}

func sumNumbersFromSlice(numberSlice []int) int64 {
	var total int64 = 0
	for _, number := range numberSlice {
		total += int64(number)
	}

	return total
}
