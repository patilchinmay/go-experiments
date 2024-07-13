package main

import "fmt"

func main() {
	resultCh := generateInts()

	// for-range automatically exits when a channel is closed.
	// No need of comma-ok syntax. It doesn't support it either.
	for i := range resultCh {
		fmt.Println(i)
	}

	fmt.Println("channel closed, exiting...")
}

func generateInts() <-chan int {
	rc := make(chan int)

	go func() {
		defer close(rc)

		for i := 0; i < 10; i++ {
			rc <- i
		}
	}()

	return rc
}
