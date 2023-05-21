package ping

import (
	"net/http"

	"github.com/go-chi/httplog"
)

type Ping struct {
}

// Ping is the handler for GET /ping
func (p *Ping) Ping(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())
	oplog.Debug().Msg("Pong")

	resp := `{"Ping":"Pong"}`

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}
