package cloudnativepatterns

import (
	"context"
	"fmt"
	"log"
	"time"
)

func Retry(cnf CloudNativeFunction, retries int, delay time.Duration) CloudNativeFunction {
	return func(ctx context.Context) error {
		for r := 0; ; r++ {
			err := cnf(ctx)

			if err == nil {
				return nil
			}

			if r >= retries {
				return fmt.Errorf("exceeded maximum number of retries: %d", retries)
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
