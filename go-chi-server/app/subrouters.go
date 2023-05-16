package app

import (
	"github.com/go-chi/chi/v5"
)

// Subrouter is implemented by all subrouters
// App exclusively works with Subrouter interface
// rather than concrete subrouters (Liskov Substitution Principle)
type Subrouter interface {
	MountOn(chi.Router)
	Getpath() string
}

// AppendSubrouter appends the subrouter to app.Subrouters
// This is useful if we have multiple subrouters
// All of them will be maintained in app.Subrouters
// so that they can be easily mounted
func (a *App) AppendSubrouter(sr Subrouter) *App {
	a.Subrouters = append(a.Subrouters, sr)

	return a
}

// MountSubrouters initializes and
// registers the subrouters onto the main router
func (a *App) MountSubrouters() {
	for _, sr := range a.Subrouters {
		sr.MountOn(a.Router)
		a.logger.Debug().Str("path", sr.Getpath()).Msg("Registered subrouter")
	}
}
