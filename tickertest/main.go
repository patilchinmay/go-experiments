package main

import (
	"context"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
)

// Define constants
const (
	tickerInterval  time.Duration = 1 * time.Second
	timeoutInterval time.Duration = 5 * time.Second
)

func main() {
	// Use the real clock for actual execution
	realClock := clock.New()

	// Create a context with a timeout
	ctxTimeout, cancel := context.WithTimeout(context.Background(), timeoutInterval)
	defer cancel()

	resultCh := ExecuteOnTick(ctxTimeout, realClock, tickerInterval)

	for {
		select {
		case result, ok := <-resultCh:
			if !ok {
				fmt.Println("result channel closed, exiting main...")
				return
			}
			fmt.Printf("received: %s\n", result)
		case <-ctxTimeout.Done():
			fmt.Println("[main] context expired, waiting for resultCh to close...")
			// Even though the main function hit timeout, the ExecuteOnTick() function may be busy in different select cases and not necessarily hit the timeout immediately until the next cycle of select executes.
			// Thus if we do not wait for the ExecuteOnTick() function to realize that a context timeout has been hit and exit properly, the main will exit and resultCh and goroutine will not close properly, resulting in possible memory leak.
			// To allow ExecuteOnTick() function to realize context timeout and properly exit,
			// we wait on resultCh to close (signifying both channel and goroutine closure due to "defer close(resultCh)" statement).
			// Drain the resultCh (as it may have unread values). Once resultCh is closed, the for-range loop will exit automatically.
			for v := range resultCh {
				fmt.Printf("[main] STILL draining resultCh.\n\tReceived: %s\n", v)
			}
			fmt.Println("[main] resultCh closed.")
			fmt.Println("[main] exiting...")
			return
		}
	}
}

func ExecuteOnTick(ctx context.Context, clock clock.Clock, tickerInterval time.Duration) <-chan string {
	// Channel used to receive the result from ExecuteOnTick function
	resultCh := make(chan string)

	// Time ticker
	ticker := clock.Ticker(tickerInterval)

	// Why execute "for-select" as goroutine?

	// It allows us to return the resultCh
	// This follows the best practice of goroutines
	// i.e. The owner of the goroutine should write and close it, and the consumer of the goroutine should only read it (and never write/close it).

	// If we do not execute "for-select" as goroutine (and put it directly in ExecuteOnTick function), then we can not return resultCh as the code flow will be blocked by "for-select" below.

	// In other words, to return the results on a channel (while "for-select" is directly in ExecuteOnTick function), we would need to accept resultCh channel as a parameter in ExecuteOnTick and write to it.
	// This will deviate from the best practices since the main function will be the owner of the channel and ExecuteOnTick will write to it.
	go func() {
		defer ticker.Stop()
		defer close(resultCh)

		for {
			select {
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				// Simulate work and sending a response
				select {
				case resultCh <- "OK":
					// Sent successfully
				case <-ctx.Done():
					fmt.Println("[function] Context expired while trying to send, exiting...")
					return
				}
			case <-ctx.Done():
				fmt.Println("[function] context expired, exiting...")
				return
			}
		}
	}()

	return resultCh
}
