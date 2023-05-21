package user

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
)

type User struct {
}

// Get is the handler for GET /user
func (u *User) Get(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())
	oplog.Debug().Msg("Add User")

	requestID := middleware.GetReqID(r.Context())
	oplog.Debug().Str("requestID", requestID).Msg("")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("requestID", requestID)

	resp := fmt.Sprintf(`{"User":"Get","requestID":%s}`, requestID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}

// Add is the handler for POST /user
func (u *User) Add(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())
	oplog.Debug().Msg("Add User")

	requestID := middleware.GetReqID(r.Context())
	oplog.Debug().Str("requestID", requestID).Msg("")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("requestID", requestID)

	resp := fmt.Sprintf(`{"User":"Add","requestID":%s}`, requestID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}
