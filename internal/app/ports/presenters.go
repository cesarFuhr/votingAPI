package ports

import (
	"encoding/json"
	"net/http"
)

// HTTPError Exception formatter to all http badRequests
type HTTPError struct {
	Message string `json:"message"`
}

// HTTPCreateAgendaReq json http representation of a create agenda request
type HTTPCreateAgendaReq struct {
	Description string `json:"description"`
}

// HTTPCreateAgendaRes json http representation of a create agenda response
type HTTPCreateAgendaRes struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func internalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(HTTPError{
		Message: "There was an unexpected error",
	})
}
