package ports

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/cesarFuhr/votingAPI/internal/app/domain/agenda"
)

type AgendaServiceStub struct {
	CalledWith          []interface{}
	LastDeliveredAgenda agenda.Agenda
}

func (s *AgendaServiceStub) CreateAgenda(description string) (agenda.Agenda, error) {
	s.CalledWith = []interface{}{description}
	if description == "ERROR" {
		return agenda.Agenda{}, errors.New("A ERROR")
	}
	return agenda.Agenda{
		ID:          "36df597d-a3b7-45cd-b65a-439c0900649e",
		Description: description,
	}, nil
}

func (s *AgendaServiceStub) FindAgenda(id string) (agenda.Agenda, error) {
	s.CalledWith = []interface{}{id}
	if id == "notFound" {
		return agenda.Agenda{}, errors.New("Agenda not found")
	}
	if id == "otherError" {
		return agenda.Agenda{}, errors.New("Any error at all")
	}

	s.LastDeliveredAgenda = agenda.Agenda{
		ID:          id,
		Description: "a description",
	}
	return s.LastDeliveredAgenda, nil
}

var validAgendaReqBody, _ = json.Marshal(HTTPCreateAgendaReq{
	Description: "a description",
})

func TestPOSTAgenda(t *testing.T) {
	agendaService := AgendaServiceStub{}
	h := NewAgendaHandler(&agendaService)
	t.Run("Should return 201 on /agenda", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda", bytes.NewBuffer(validAgendaReqBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusCreated)
	})
	t.Run("Should return valid json on /agenda", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda", bytes.NewBuffer(validAgendaReqBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		respBytes, _ := ioutil.ReadAll(response.Body)

		if !json.Valid(respBytes) {
			t.Errorf("got an invalid JSON %q", respBytes)
		}
	})
	t.Run("Should return all properties on /agenda response", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda", bytes.NewBuffer(validAgendaReqBody))
		response := httptest.NewRecorder()

		wants := []string{"id", "description"}

		h.Post(response, request)
		respMap := map[string]interface{}{}
		extractJSON(response.Body, respMap)

		for _, want := range wants {
			if _, ok := respMap[want]; ok != true {
				t.Errorf("does not have the %q prop", want)
			}
		}
	})
	t.Run("Should call the CreateAgenda with the correct params", func(t *testing.T) {
		description := "a description"
		requestBody, _ := json.Marshal(map[string]string{
			"description": description,
		})
		request, _ := http.NewRequest(http.MethodPost, "/agenda", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertInsideSlice(t, agendaService.CalledWith, description)
	})
	t.Run("Should return a internal server error if there was an error creating an agenda", func(t *testing.T) {
		description := "ERROR"
		requestBody, _ := json.Marshal(map[string]string{
			"description": description,
		})
		request, _ := http.NewRequest(http.MethodPost, "/agenda", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
	})
	t.Run("Should return a BadRequest if body is nil", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda", nil)
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertInsideJSON(t, response.Body, "message", "Invalid: Empty body")
	})
}

func TestGETAgenda(t *testing.T) {
	agendaService := AgendaServiceStub{}
	h := NewAgendaHandler(&agendaService)
	t.Run("Should return a 200 if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/agenda/anID", nil)
		response := httptest.NewRecorder()

		want := http.StatusOK

		h.Get(response, request)

		assertStatus(t, response.Code, want)
	})
	t.Run("Should return a Agenda if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/agenda/anID", nil)
		response := httptest.NewRecorder()

		wants := []string{"id", "description"}

		h.Get(response, request)
		respMap := map[string]interface{}{}
		extractJSON(response.Body, respMap)

		for _, want := range wants {
			if _, ok := respMap[want]; ok != true {
				t.Errorf("does not have the %q prop", want)
			}
		}
	})
	t.Run("Should call find Agenda with the right params", func(t *testing.T) {
		getRequest, _ := http.NewRequest(http.MethodGet, "/agenda/anID", nil)
		response := httptest.NewRecorder()
		h.Get(response, getRequest)

		assertInsideSlice(t, agendaService.CalledWith, "anID")
	})
	t.Run("If agenda was not found", func(t *testing.T) {
		t.Run("Should return a 404", func(t *testing.T) {
			want := http.StatusNotFound
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/agenda/notFound"), nil)
			response := httptest.NewRecorder()
			h.Get(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "Agenda not found")
		})
	})
	t.Run("If there was any other error", func(t *testing.T) {
		t.Run("Should return a 500", func(t *testing.T) {
			want := http.StatusInternalServerError
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/agenda/otherError"), nil)
			response := httptest.NewRecorder()
			h.Get(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
		})
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}

func extractJSON(jBuff *bytes.Buffer, m map[string]interface{}) error {
	respBytes, _ := ioutil.ReadAll(jBuff)
	_ = json.Unmarshal(respBytes, &m)
	return nil
}

func assertInsideJSON(t *testing.T, jBuff *bytes.Buffer, wantedKey string, wantedValue interface{}) {
	t.Helper()
	got := map[string]interface{}{}
	extractJSON(jBuff, got)
	if !reflect.DeepEqual(got[wantedKey], wantedValue) {
		t.Errorf("got %v, want %v", got[wantedKey], wantedValue)
	}
}

func assertInsideSlice(t *testing.T, a []interface{}, want interface{}) {
	t.Helper()
	has := false
	for _, v := range a {
		if v == want {
			has = true
		}
	}
	if !has {
		t.Errorf("Did not found: %v, of type %T in %v", want, want, a)
	}
}
