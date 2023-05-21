package user

import (
	"net/http"

	"github.com/go-chi/httplog"
)

type User struct {
}

// Get is the handler for GET /user
func (u *User) Get(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())
	oplog.Debug().Msg("Get User")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := `{"User":"Get"}`
	w.Write([]byte(resp))
}

// Add is the handler for POST /user
func (u *User) Add(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())
	oplog.Debug().Msg("Add User")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := `{"User":"Add"}`
	w.Write([]byte(resp))
}
