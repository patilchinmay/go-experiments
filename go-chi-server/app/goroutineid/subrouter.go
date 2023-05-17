package goroutineid

import (
	"github.com/go-chi/chi/v5"
	"github.com/patilchinmay/go-experiments/go-chi-server/app"
)

// Goroutinecheck implements app.Subrouter interface
type Goroutinecheck struct {
	Path      string
	Subrouter chi.Router
}

var goroutinecheck *Goroutinecheck

// GetOrCreate returns a singleton instance of Ping
func GetOrCreate() *Goroutinecheck {
	if goroutinecheck == nil {
		goroutinecheck = &Goroutinecheck{
			Path:      "/goroutinecheck",
			Subrouter: chi.NewRouter(),
		}
	}
	return goroutinecheck
}

// InitializeRoutes associates the http.HandlerFuncs
// from Ping package and http.Method with Ping.Subrouter
func (g *Goroutinecheck) InitializeRoutes() *Goroutinecheck {
	g.Subrouter.Get("/", g.CheckGoroutineID) // GET /ping

	return g
}

// MountOn mounts the Ping.Subrouter onto r (the main router)
func (g *Goroutinecheck) MountOn(r chi.Router) {
	r.Mount(g.Path, g.Subrouter) // r is the main router
}

// Getpath returns the Ping.Path
func (g *Goroutinecheck) Getpath() string {
	return g.Path
}

// init initializes the ping subrouter and
// appends it to the []app.Subrouters as side-effects.
// This function will be executed automatically
// when this package is imported (as a dash import/blank identifier).
func init() {
	// Create ping subrouter with routes
	goroutinecheck := GetOrCreate().InitializeRoutes()
	// Append the ping subrouter to []app.Subrouters
	app.GetOrCreate().AppendSubrouter(goroutinecheck)
}
