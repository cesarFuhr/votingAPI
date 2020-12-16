package ports

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/vote"
)

type VoteServiceStub struct {
	CalledWith        []interface{}
	LastDeliveredVote vote.Vote
}

func (s *VoteServiceStub) CreateVote(associateID, sessionID, document, value string) (vote.Vote, error) {
	s.CalledWith = []interface{}{associateID, sessionID, document, value}
	if sessionID == "ERROR" {
		return vote.Vote{}, errors.New("A ERROR")
	}
	if sessionID == "notFound" {
		return vote.Vote{}, errors.New("Session not found")
	}
	if sessionID == "duplicateVote" {
		return vote.Vote{}, vote.ErrDuplicateVote
	}
	return vote.Vote{
		AssociateID: associateID,
		SessionID:   sessionID,
		Document:    document,
		Vote:        value,
		Creation:    time.Now(),
	}, nil
}

var validVoteReqBody, _ = json.Marshal(HTTPCreateVoteReq{
	AssociateID: "associateID",
	Document:    "01212393111",
	Vote:        "S",
})

func TestPOSTVote(t *testing.T) {
	voteService := VoteServiceStub{}
	h := NewVoteHandler(&voteService)
	t.Run("Should return 201 on /agenda/id/session/id/vote", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda/id/session/id/vote", bytes.NewBuffer(validVoteReqBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusCreated)
	})
	t.Run("Should call the CreateVote with the correct params", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "agenda/id/session/sessionID/vote", bytes.NewBuffer(validVoteReqBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertInsideSlice(t, voteService.CalledWith, "associateID")
		assertInsideSlice(t, voteService.CalledWith, "01212393111")
		assertInsideSlice(t, voteService.CalledWith, "sessionID")
		assertInsideSlice(t, voteService.CalledWith, "S")
	})
	t.Run("Should return a internal server error if there was an error creating an vote", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda/id/session/ERROR/vote", bytes.NewBuffer(validVoteReqBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
	})
	t.Run("Should return a internal server error if there was an error creating an vote", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda/id/session/notFound/vote", bytes.NewBuffer(validVoteReqBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertInsideJSON(t, response.Body, "message", "Session not found")
	})
	t.Run("Should return a internal server error if there was an error creating an vote", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda/id/session/duplicateVote/vote", bytes.NewBuffer(validVoteReqBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertInsideJSON(t, response.Body, "message", vote.ErrDuplicateVote.Error())
	})
}
