package ports

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/session"
)

func (s *SessionServiceStub) Result(id string) (session.Result, error) {
	s.CalledWith = []interface{}{id}
	if id == "notFound" {
		return session.Result{}, errors.New("Session not found")
	}
	if id == "otherError" {
		return session.Result{}, errors.New("Any error at all")
	}

	return session.Result{
		ID:             "36df597d-a3b7-45cd-b65a-439c0900649e",
		OriginalAgenda: "originalAgenda",
		Closed:         true,
		Count:          session.Count{InFavor: 10, Against: 12},
	}, nil
}

func TestGETResult(t *testing.T) {
	sessionService := SessionServiceStub{}
	h := NewResultHandler(&sessionService)
	t.Run("Should return a 200 if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/agenda/id/session/anID/result", nil)
		response := httptest.NewRecorder()

		want := http.StatusOK

		h.Get(response, request)

		assertStatus(t, response.Code, want)
	})
	t.Run("Should return a Result if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/agenda/id/session/anID/result", nil)
		response := httptest.NewRecorder()

		wants := []string{"id", "originalAgenda", "closed", "count"}

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
		getRequest, _ := http.NewRequest(http.MethodGet, "/agenda/anotherID/session/anID/result", nil)
		response := httptest.NewRecorder()
		h.Get(response, getRequest)

		assertInsideSlice(t, sessionService.CalledWith, "anID")
	})
	t.Run("If session was not found", func(t *testing.T) {
		t.Run("Should return a 404", func(t *testing.T) {
			want := http.StatusNotFound
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/agenda/id/session/notFound/result"), nil)
			response := httptest.NewRecorder()
			h.Get(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "Session not found")
		})
	})
	t.Run("If there was any other error", func(t *testing.T) {
		t.Run("Should return a 500", func(t *testing.T) {
			want := http.StatusInternalServerError
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/agenda/id/session/otherError/result"), nil)
			response := httptest.NewRecorder()
			h.Get(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
		})
	})
}
