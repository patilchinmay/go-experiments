package app

import (
	"github.com/go-chi/chi/v5"
)

// Subrouter definition
type Subrouter struct {
	Path      string
	Subrouter chi.Router
}

// NewSubrouter is a constructor for Subrouter
func NewSubrouter(mountpath string) Subrouter {
	var sr = Subrouter{
		Path:      mountpath,
		Subrouter: chi.NewRouter(),
	}

	return sr
}

// AppendSubrouter appends the subrouter to app.Subrouters
// This is useful if we have multiple subrouters
// All of them will be maintained in app.Subrouters
// so that they can be easily mounted
func (a *App) AppendSubrouter(sr Subrouter) *App {
	a.Subrouters = append(a.Subrouters, sr)

	return a
}

// MountSubrouters initializes and registers
// the subrouters onto the main router.
// Middlewares must be defined before this function is called.
func (a *App) MountSubrouters() {
	for _, sr := range a.Subrouters {
		a.Router.Mount(sr.Path, sr.Subrouter)
		a.logger.Debug().Str("path", sr.Path).Msg("Registered subrouter")
	}
}
