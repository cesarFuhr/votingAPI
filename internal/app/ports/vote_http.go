package ports

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cesarFuhr/votingAPI/internal/app/domain/vote"
)

type voteHandler struct {
	service vote.Service
}

// VoteHandler describes a http handler interface
type VoteHandler interface {
	Post(w http.ResponseWriter, r *http.Request)
}

// NewVoteHandler creates a new http vote handler
func NewVoteHandler(s vote.Service) VoteHandler {
	return &voteHandler{
		service: s,
	}
}

// Post http translator
func (h *voteHandler) Post(w http.ResponseWriter, r *http.Request) {
	trimmed := strings.SplitAfter(r.URL.Path, "/session/")
	sessionID := strings.TrimSuffix(trimmed[1], "/vote")

	var o HTTPCreateVoteReq
	err := decodeJSONBody(r, &o, true)
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

	_, err = h.service.CreateVote(o.AssociateID, sessionID, o.Document, o.Vote)
	if err != nil {
		if err == vote.ErrDuplicateVote ||
			err == vote.ErrBadVoteFormat ||
			err == vote.ErrSessionExpired ||
			err == vote.ErrNotAbleToVote {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(HTTPError{
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "Session not found" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(HTTPError{
				Message: err.Error(),
			})
			return
		}
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
