package user

import (
	"github.com/patilchinmay/go-experiments/go-chi-server/app"
)

// init initializes the subrouter, defines the routes & handlers, and
// appends it to the []app.Subrouters.
// This function will be executed automatically
// when this package is imported (as a dash import/blank identifier).
func init() {
	path := "/user"

	// Create subrouter with routes
	sr := app.NewSubrouter(path)

	// Register methods
	u := Users{}

	// Define the routes on subrouter
	// All the routes here have a prefix of
	// path defined above.
	sr.Subrouter.Get("/", u.Get)
	sr.Subrouter.Post("/", u.Add)

	// Append to app
	app.GetOrCreate().AppendSubrouter(sr)
}
