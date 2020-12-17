package ports

import (
	"encoding/json"
	"net/http"
	"time"
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

// HTTPCreateSessionReq json http representation of a create session request
type HTTPCreateSessionReq struct {
	Duration time.Duration `json:"durationInMinutes"`
}

// HTTPCreateSessionRes json http representation of a create session response
type HTTPCreateSessionRes struct {
	ID             string `json:"id"`
	OriginalAgenda string `json:"originalAgenda"`
	Expiration     string `json:"expiration"`
}

// HTTPCreateVoteReq json http representation of a create vote request
type HTTPCreateVoteReq struct {
	AssociateID string `json:"associateID"`
	Document    string `json:"document"`
	Vote        string `json:"vote"`
}

// HTTPResultSessionRes json http representation of a session result response
type HTTPResultSessionRes struct {
	ID             string `json:"id"`
	OriginalAgenda string `json:"originalAgenda"`
	Closed         bool   `json:"closed"`
	Count          struct {
		InFavor  int `json:"inFavor"`
		Angainst int `json:"against"`
	} `json:"count"`
}

func internalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(HTTPError{
		Message: "There was an unexpected error",
	})
}
