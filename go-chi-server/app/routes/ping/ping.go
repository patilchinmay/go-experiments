package ping

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/patilchinmay/go-experiments/go-chi-server/app"
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

// InitializeRoutes adds the routes
// from current files in local subrouter
func (p *Ping) InitializeRoutes() *Ping {
	p.Subrouter.Get(p.Path, p.Ping) // GET /ping

	return p
}

// MountSubrouter mounts the subrouter onto the main router
func (p *Ping) MountSubrouterOn(r chi.Router) {
	r.Mount("/", p.Subrouter) // r is the main router
}

// AddToAppSubrouters appends this Subrouter to app.Subrouters
func (p *Ping) AddToAppSubrouters(a *app.App) {
	a.Subrouters = append(a.Subrouters, p)
}

func (p *Ping) Getpath() string {
	return p.Path
}

func (p *Ping) Ping(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())
	oplog.Debug().Msg("Pong")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte("Pong"))
}
