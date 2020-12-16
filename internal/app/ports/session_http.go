package ports

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/session"
)

type sessionOpts struct {
	Duration time.Duration `json:"durationInMinutes"`
}

type sessionHandler struct {
	service session.Service
}

// SessionHandler describes a http handler interface
type SessionHandler interface {
	Post(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

// NewSessionHandler creates a new http session handler
func NewSessionHandler(s session.Service) SessionHandler {
	return &sessionHandler{
		service: s,
	}
}

// Post http translator
func (h *sessionHandler) Post(w http.ResponseWriter, r *http.Request) {
	trimmed := strings.SplitAfter(r.URL.Path, "/agenda/")
	originalAgenda := strings.TrimSuffix(trimmed[1], "/session")

	var o sessionOpts
	err := decodeJSONBody(r, &o, true)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			w.WriteHeader(mr.status)
			json.NewEncoder(w).Encode(HTTPError{
				Message: mr.msg,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HTTPError{
			Message: fmt.Sprint(err),
		})
		return
	}

	session, err := h.service.CreateSession(originalAgenda, o.Duration)
	if err != nil {
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(HTTPCreateSessionRes{
		ID:             session.ID,
		OriginalAgenda: session.OriginalAgenda,
		Expiration:     session.GetExpiration().Format(time.RFC3339),
	})
	return
}

func (h *sessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	sliced := strings.Split(r.URL.Path, "/session/")
	id := sliced[1]

	session, err := h.service.FindSession(id)
	if err != nil {
		if err.Error() == "Session not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(HTTPError{
				Message: "Session not found",
			})
			return
		}
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HTTPCreateSessionRes{
		ID:             session.ID,
		OriginalAgenda: session.OriginalAgenda,
		Expiration:     session.GetExpiration().Format(time.RFC3339),
	})
	return
}
