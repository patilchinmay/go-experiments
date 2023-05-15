package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/patilchinmay/go-experiments/go-chi-server/app/routes/ping"
)

// Subrouter is implemented by all subrouters
// App exclusively works with Subrouter interface
// rather than concrete subrouters (Liskov Substitution Principle)
type Subrouter interface {
	MountOn(chi.Router)
	Getpath() string
}

// CreateSubrouters creates an instance of the subrouters
func (a *App) CreateSubrouters() *App {
	// Create ping subrouter
	// Register path and handler
	// Add the subrouter to app.Subrouters
	ping := ping.New().InitializeRoutes()
	a.AppendSubrouter(ping)

	// Any additional subrouter should be initialized here
	return a
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
