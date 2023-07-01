package user

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

type UserHandler struct {
	usrsvc *UserService
}

var userHandler *UserHandler

func NewUserHandler(usrsvc *UserService) *UserHandler {
	if userHandler == nil {
		userHandler = &UserHandler{
			usrsvc: usrsvc,
		}
	}
	return userHandler
}

// Get is the handler for GET /user/{id}
func (u *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())
	oplog.Debug().Msg("Get User")

	// Retrieve the path param of id. id is string here.
	idstr := chi.URLParam(r, "id")
	if idstr == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		oplog.Error().Msg("Missing id")
		return
	}

	// convert id to uint (as expected by the service layer)
	id64, err := strconv.ParseUint(idstr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		oplog.Error().Msg("Invalid id")
		return
	}
	id := uint(id64)

	// Call the service layer and retrieve the user
	user, err := u.usrsvc.get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		oplog.Error().Err(err).Msg("Failed to get user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Add is the handler for POST /user
func (u *UserHandler) Add(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())
	oplog.Debug().Msg("Add User")

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		oplog.Error().Err(err).Msg("Failed to read request body")
		return
	}

	// Parse the incoming payload into json
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		// Parse []byte to go struct pointer
		http.Error(w, err.Error(), http.StatusBadRequest)
		oplog.Error().Err(err).Msg("Failed to parse body as json")
	}

	w.Header().Set("Content-Type", "application/json")

	// Call the service layer
	resp, err := u.usrsvc.add(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		oplog.Error().Err(err).Msg("Body failed validation")
		return
	}

	// Return the response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.FormatUint(uint64(resp), 10)))
}
