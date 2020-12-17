package ports

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/agenda"
)

type agendaOpts struct {
	Description string `json:"description"`
}

type agendaHandler struct {
	service agenda.Service
}

// AgendaHandler describes a http handler interface
type AgendaHandler interface {
	Post(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

// NewAgendaHandler creates a new http agenda handler
func NewAgendaHandler(s agenda.Service) AgendaHandler {
	return &agendaHandler{
		service: s,
	}
}

// Post http translator
func (h *agendaHandler) Post(w http.ResponseWriter, r *http.Request) {
	var o agendaOpts
	err := decodeJSONBody(r, &o, false)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(mr.status)
			json.NewEncoder(w).Encode(HTTPError{
				Message: mr.msg,
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HTTPError{
			Message: fmt.Sprint(err),
		})
		return
	}

	agenda, err := h.service.CreateAgenda(o.Description)
	if err != nil {
		internalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(HTTPCreateAgendaRes{
		ID:          agenda.ID,
		Description: agenda.Description,
	})
	return
}

func (h *agendaHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/agenda/")

	agenda, err := h.service.FindAgenda(id)
	if err != nil {
		if err.Error() == "Agenda not found" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(HTTPError{
				Message: "Agenda not found",
			})
			return
		}
		internalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HTTPCreateAgendaRes{
		ID:          agenda.ID,
		Description: agenda.Description,
	})
	return
}
