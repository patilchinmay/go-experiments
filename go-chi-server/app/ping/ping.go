package ping

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/patilchinmay/go-experiments/go-chi-server/app"
)

// Ping implements app.Subrouter interface
type Ping struct {
	Path      string
	Subrouter chi.Router
}

var ping *Ping

// GetOrCreate returns a singleton instance of Ping
func GetOrCreate() *Ping {
	if ping == nil {
		ping = &Ping{
			Path:      "/ping",
			Subrouter: chi.NewRouter(),
		}
	}
	return ping
}

// InitializeRoutes associates the http.HandlerFuncs
// from Ping package and http.Method with Ping.Subrouter
func (p *Ping) InitializeRoutes() *Ping {
	p.Subrouter.Get(p.Path, p.Ping) // GET /ping

	return p
}

// MountOn mounts the Ping.Subrouter onto r (the main router)
func (p *Ping) MountOn(r chi.Router) {
	r.Mount("/", p.Subrouter) // r is the main router
}

// Getpath returns the Ping.Path
func (p *Ping) Getpath() string {
	return p.Path
}

// Ping is the handler for GET /ping
func (p *Ping) Ping(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())
	oplog.Debug().Msg("Pong")

	requestID := middleware.GetReqID(r.Context())
	oplog.Debug().Str("requestID", requestID).Msg("")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("requestID", requestID)

	resp := fmt.Sprintf(`{"Ping":"Pong","requestID":%s}`, requestID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}

// init initializes the ping subrouter and
// appends it to the []app.Subrouters as side-effects.
// This function will be executed automatically
// when this package is imported (as a dash import/blank identifier).
func init() {
	// Create ping subrouter with routes
	ping := GetOrCreate().InitializeRoutes()
	// Append the ping subrouter to []app.Subrouters
	app.GetOrCreate().AppendSubrouter(ping)
}
