package cloudnativepatterns

import (
	"context"
	"fmt"
	"time"
)

func (cnp *CNP) Retry(cnf CloudNativeFunction, retries int, delay time.Duration) CloudNativeFunction {
	return func(ctx context.Context) error {
		for r := 0; ; r++ {
			err := cnf(ctx)

			if err == nil {
				return nil
			}

			if r >= retries {
				return fmt.Errorf("exceeded maximum number of retries: %d", retries)
			}

			fmt.Printf("Attempt %d failed at %v; retrying in %v\n", r+1, cnp.Clock.Now(), delay)

			select {
			case <-cnp.Clock.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
