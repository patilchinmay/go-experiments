package user

import (
	"github.com/patilchinmay/go-experiments/go-chi-server/app"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// SetupSubrouter initializes the subrouter, defines the routes & handlers, and
// appends it to the []app.Subrouters.
// This function is called in main
func SetupSubrouter(db *gorm.DB, logger zerolog.Logger) {
	path := "/user"

	// Create subrouter with routes
	sr := app.NewSubrouter(path)

	// Initiate User Repository Layer
	usrrepo := NewUserRepository(db)

	// Initiate User Service
	usrsvc := NewUserService(usrrepo)

	// Initiate User handler
	usrhandler := NewUserHandler(usrsvc)

	// Define the routes on subrouter
	// All the routes here have a prefix of
	// path defined above.
	sr.Subrouter.Get("/", usrhandler.Get)
	sr.Subrouter.Post("/", usrhandler.Add)

	// Append to app
	app.GetOrCreate().AppendSubrouter(sr)
}