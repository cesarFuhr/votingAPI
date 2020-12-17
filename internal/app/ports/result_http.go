package ports

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/session"
)

type resultHandler struct {
	service session.Service
}

// ResultHandler describes a http handler interface
type ResultHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
}

// NewResultHandler creates a new http session handler
func NewResultHandler(s session.Service) ResultHandler {
	return &resultHandler{
		service: s,
	}
}

func (h *resultHandler) Get(w http.ResponseWriter, r *http.Request) {
	sliced := strings.Split(r.URL.Path, "/session/")
	id := strings.TrimSuffix(sliced[1], "/result")

	result, err := h.service.Result(id)
	if err != nil {
		if err.Error() == "Session not found" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(HTTPError{
				Message: err.Error(),
			})
			return
		}
		internalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseBody := HTTPResultSessionRes{
		ID:             result.ID,
		OriginalAgenda: result.OriginalAgenda,
		Closed:         result.Closed,
	}
	responseBody.Count.InFavor = result.Count.InFavor
	responseBody.Count.Angainst = result.Count.Against
	json.NewEncoder(w).Encode(&responseBody)
	return
}
