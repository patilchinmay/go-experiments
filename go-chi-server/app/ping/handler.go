package ping

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
)

type Ping struct {
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
