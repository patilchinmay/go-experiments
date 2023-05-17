package goroutineid

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-chi/httplog"
)

// goid returns the id of the current goroutine
func (g *Goroutinecheck) goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

// Ping is the handler for GET /ping
func (g *Goroutinecheck) CheckGoroutineID(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())

	goroutineID := g.goid()

	oplog.Debug().Int("goroutineID", goroutineID).Msg("")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("requestID", strconv.Itoa(goroutineID))

	resp := fmt.Sprintf(`{"goroutineID":%d}`, goroutineID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}
