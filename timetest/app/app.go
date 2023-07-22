package app

import (
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
)

type Application struct {
	Clock clock.Clock
}

var app *Application

func NewApplication(clock clock.Clock) *Application {
	if app == nil {
		app = &Application{
			Clock: clock,
		}
	}
	return app
}

func (a Application) Now() {
	fmt.Println(a.Clock.Now())
}

func (a Application) After(d time.Duration) {
	fmt.Println(<-a.Clock.After(d))
}
