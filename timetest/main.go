package main

import (
	"time"

	"github.com/benbjohnson/clock"
	"github.com/patilchinmay/go-experiments/timetest/app"
)

func main() {
	// for real time
	// clock := clock.New()

	// for mock time
	clock := clock.NewMock()

	app := app.NewApplication(clock)

	app.Now()

	doneCh := make(chan struct{})
	defer close(doneCh)

	// Run a separate goroutine to increase the time
	go func() {
		for {
			select {
			case <-doneCh:
				return
			default:
				clock.Add(1 * time.Second)
			}
		}
	}()

	app.After(5 * time.Second)
}
