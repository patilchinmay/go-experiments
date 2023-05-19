package validator

import (
	"github.com/go-chi/chi/v5"
	"github.com/patilchinmay/go-experiments/go-chi-server/app"
)

// Subrouter implements app.Subrouter interface
type Subrouter struct {
	Path      string
	Subrouter chi.Router
}

var subrouter *Subrouter

// GetOrCreate returns a singleton instance of Subrouter
func GetOrCreate() *Subrouter {
	if subrouter == nil {
		subrouter = &Subrouter{
			Path:      "/validator",
			Subrouter: chi.NewRouter(),
		}
	}
	return subrouter
}

// InitializeRoutes associates the http.HandlerFuncs
// from Subrouter package and http.Method with Subrouter.Subrouter
func (sr *Subrouter) InitializeRoutes() *Subrouter {
	v := Validator{}
	sr.Subrouter.Get("/", v.Validate)

	return sr
}

// MountOn mounts the Subrouter.Subrouter onto r (the main router)
func (sr *Subrouter) MountOn(r chi.Router) {
	r.Mount(sr.Path, sr.Subrouter) // r is the main router
}

// Getpath returns the Subrouter.Path
func (sr *Subrouter) Getpath() string {
	return sr.Path
}

// init initializes the subrouter and
// appends it to the []app.Subrouters as side-effects.
// This function will be executed automatically
// when this package is imported (as a dash import/blank identifier).
func init() {
	// Create subrouter with routes
	subrouter := GetOrCreate().InitializeRoutes()
	// Append the subrouter subrouter to []app.Subrouters
	app.GetOrCreate().AppendSubrouter(subrouter)
}
