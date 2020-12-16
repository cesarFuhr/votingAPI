package ports

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/session"
)

type SessionServiceStub struct {
	CalledWith           []interface{}
	LastDeliveredSession session.Session
}

func (s *SessionServiceStub) CreateSession(originalAgenda string, duration time.Duration) (session.Session, error) {
	s.CalledWith = []interface{}{originalAgenda, duration}
	if originalAgenda == "ERROR" {
		return session.Session{}, errors.New("A ERROR")
	}
	return session.Session{
		ID:             "36df597d-a3b7-45cd-b65a-439c0900649e",
		OriginalAgenda: originalAgenda,
		Creation:       time.Now(),
		Duration:       time.Minute,
	}, nil
}

func (s *SessionServiceStub) FindSession(id string) (session.Session, error) {
	s.CalledWith = []interface{}{id}
	if id == "notFound" {
		return session.Session{}, errors.New("Session not found")
	}
	if id == "otherError" {
		return session.Session{}, errors.New("Any error at all")
	}

	s.LastDeliveredSession = session.Session{
		ID:             "36df597d-a3b7-45cd-b65a-439c0900649e",
		OriginalAgenda: "originalAgenda",
		Creation:       time.Now(),
		Duration:       time.Minute,
	}
	return s.LastDeliveredSession, nil
}

var validSessionReqBody, _ = json.Marshal(HTTPCreateSessionReq{
	Duration: time.Minute,
})

func TestPOSTSession(t *testing.T) {
	sessionService := SessionServiceStub{}
	h := NewSessionHandler(&sessionService)
	t.Run("Should return 201 on /session", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(validSessionReqBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusCreated)
	})
	t.Run("Should return valid json on /agenda/:id/session", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda/id/session", bytes.NewBuffer(validSessionReqBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		respBytes, _ := ioutil.ReadAll(response.Body)

		if !json.Valid(respBytes) {
			t.Errorf("got an invalid JSON %q", respBytes)
		}
	})
	t.Run("Should return all properties on /session response", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda/id/session", bytes.NewBuffer(validSessionReqBody))
		response := httptest.NewRecorder()

		wants := []string{"id", "originalAgenda", "expiration"}

		h.Post(response, request)
		respMap := map[string]interface{}{}
		extractJSON(response.Body, respMap)

		for _, want := range wants {
			if _, ok := respMap[want]; ok != true {
				t.Errorf("does not have the %q prop", want)
			}
		}
	})
	t.Run("Should call the CreateSession with expiration and scope", func(t *testing.T) {
		duration := time.Minute
		requestBody, _ := json.Marshal(map[string]interface{}{
			"durationInMinutes": duration,
		})
		request, _ := http.NewRequest(http.MethodPost, "agenda/id/session", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertInsideSlice(t, sessionService.CalledWith, duration)
	})
	t.Run("Should return a internal server error if there was an error creating an session", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda/ERROR/session", nil)
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
	})
}

func TestGETSession(t *testing.T) {
	sessionService := SessionServiceStub{}
	h := NewSessionHandler(&sessionService)
	t.Run("Should return a 200 if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/session/anID", nil)
		response := httptest.NewRecorder()

		want := http.StatusOK

		h.Get(response, request)

		assertStatus(t, response.Code, want)
	})
	t.Run("Should return a Session if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/session/anID", nil)
		response := httptest.NewRecorder()

		wants := []string{"id", "originalAgenda", "expiration"}

		h.Get(response, request)
		respMap := map[string]interface{}{}
		extractJSON(response.Body, respMap)

		for _, want := range wants {
			if _, ok := respMap[want]; ok != true {
				t.Errorf("does not have the %q prop", want)
			}
		}
	})
	t.Run("Should call find Session with the right params", func(t *testing.T) {
		getRequest, _ := http.NewRequest(http.MethodGet, "/agenda/anotherID/session/anID", nil)
		response := httptest.NewRecorder()
		h.Get(response, getRequest)

		assertInsideSlice(t, sessionService.CalledWith, "anID")
	})
	t.Run("If session was not found", func(t *testing.T) {
		t.Run("Should return a 404", func(t *testing.T) {
			want := http.StatusNotFound
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/agenda/id/session/notFound"), nil)
			response := httptest.NewRecorder()
			h.Get(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "Session not found")
		})
	})
	t.Run("If there was any other error", func(t *testing.T) {
		t.Run("Should return a 500", func(t *testing.T) {
			want := http.StatusInternalServerError
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/agenda/id/session/otherError"), nil)
			response := httptest.NewRecorder()
			h.Get(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
		})
	})
}
