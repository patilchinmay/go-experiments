package validator

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httplog"
)

// Validate is the handler for GET /validator
func (g *Validator) Validate(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())

	requestID := middleware.GetReqID(r.Context())
	w.Header().Set("requestID", requestID)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	// Write JSON response
}
