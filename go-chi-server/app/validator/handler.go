package validator

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate // TODO: Move this into main app
}

// Validate is the handler for GET /validator
func (v *Validator) Validate(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())

	requestID := middleware.GetReqID(r.Context())
	w.Header().Set("requestID", requestID)

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
		oplog.Error().Err(err).Msg("Failed to parse body as json")
	}

	// Validate the json
	v.validate = validator.New()

	// returns nil or ValidationErrors ( []FieldError )
	err = v.validate.Struct(user)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			oplog.Error().Err(err).Msg("Failed validation")
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			oplog.Error().Err(err).Msg(err.Error())

			// fmt.Println(err.Namespace())
			// fmt.Println(err.Field())
			// fmt.Println(err.StructNamespace())
			// fmt.Println(err.StructField())
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())
			// fmt.Println()
		}

		// from here you can create your own error messages in whatever language you wish
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
