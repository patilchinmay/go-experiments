package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// RequestID adds the requestID header to each outgoing response
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// oplog := httplog.LogEntry(r.Context())

		// oplog.Debug().Msg("Executed before request is passed to handler")

		requestID := middleware.GetReqID(r.Context())
		w.Header().Set("Request-ID", requestID)

		next.ServeHTTP(w, r)

		// oplog.Debug().Msg("Executed after request is passed to the handler")
	})
}
