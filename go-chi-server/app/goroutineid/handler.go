package goroutineid

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-chi/httplog"
)

type Goroutineid struct {
}

// goid returns the id of the current goroutine
func (g *Goroutineid) goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

// CheckGoroutineID is the handler for GET /goroutinecheck
func (g *Goroutineid) CheckGoroutineID(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())

	goroutineID := g.goid()

	oplog.Debug().Int("goroutineID", goroutineID).Msg("")

	w.Header().Set("Content-Type", "application/json")

	resp := fmt.Sprintf(`{"goroutineID":%d}`, goroutineID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}
