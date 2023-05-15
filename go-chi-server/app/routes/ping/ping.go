package ping

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

// Ping implements app.Subrouter interface
type Ping struct {
	Path      string
	Subrouter chi.Router
}

var ping *Ping

// New returns a singleton instance of Ping
func New() *Ping {
	if ping == nil {
		return &Ping{
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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte("Pong"))
}
