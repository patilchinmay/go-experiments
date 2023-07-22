package cloudnativepatterns

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
)

type CloudNativeFunction func(context.Context) error

type CloudNativePatterns interface {
	Retry(cnf CloudNativeFunction, retries int, delay time.Duration) CloudNativeFunction
}

type CNP struct {
	Clock clock.Clock
}

var cnp *CNP

func NewCloudNativePatterns(clock clock.Clock) *CNP {
	if cnp == nil {
		cnp = &CNP{
			Clock: clock,
		}
	}
	return cnp
}

func DiscardCloudNativePatterns() {
	if cnp != nil {
		cnp = nil
	}
}
