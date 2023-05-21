package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// RequestID adds the requestID header to each outgoing response
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// oplog := httplog.LogEntry(r.Context())

		requestID := middleware.GetReqID(r.Context())
		// oplog.Debug().Str("requestID", requestID).Msg("")

		w.Header().Set("requestID", requestID)
		next.ServeHTTP(w, r)
	})
}
