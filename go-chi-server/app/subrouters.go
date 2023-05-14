package app

import (
	"github.com/go-chi/chi/v5"
)

type Subrouter interface {
	MountSubrouterOn(chi.Router)
	AddToAppSubrouters(*App)
	Getpath() string
}

// RegisterSubroutersOn initializes and
// registers the subrouters onto the main router
func (a *App) RegisterSubrouters() {
	for _, sr := range a.Subrouters {
		sr.MountSubrouterOn(a.Router)
		a.logger.Debug().Str("path", sr.Getpath()).Msg("Registered subrouter")
	}
}
