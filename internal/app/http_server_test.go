package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type agendaHandlerStub struct {
	G struct {
		CalledWith []interface{}
	}
	P struct {
		CalledWith []interface{}
	}
}

func (h *agendaHandlerStub) Post(w http.ResponseWriter, r *http.Request) {
	h.P.CalledWith = []interface{}{w, r}
}

func (h *agendaHandlerStub) Get(w http.ResponseWriter, r *http.Request) {
	h.G.CalledWith = []interface{}{w, r}
}

type sessionHandlerStub struct {
	G struct {
		CalledWith []interface{}
	}
	P struct {
		CalledWith []interface{}
	}
}

func (h *sessionHandlerStub) Post(w http.ResponseWriter, r *http.Request) {
	h.P.CalledWith = []interface{}{w, r}
}

func (h *sessionHandlerStub) Get(w http.ResponseWriter, r *http.Request) {
	h.G.CalledWith = []interface{}{w, r}
}

type loggerStub struct {
	CalledWith []interface{}
}

func (l *loggerStub) Info(args ...interface{}) {
	l.CalledWith = args
}

var (
	log = loggerStub{}
	aH  = agendaHandlerStub{}
	sH  = sessionHandlerStub{}
)

func TestAgendaEndpoint(t *testing.T) {
	server := NewHTTPServer(&log, &aH, &sH)
	t.Run("calls agendaHandler.Post in a /agenda http POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertInsideSlice(t, aH.P.CalledWith, response)
		assertInsideSlice(t, aH.P.CalledWith, request)
	})
	t.Run("calls agendaHandler.Get in a /agenda http GET", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/agenda/anID", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertInsideSlice(t, aH.G.CalledWith, response)
		assertInsideSlice(t, aH.G.CalledWith, request)
	})
	t.Run("returns method not allowed for any other method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/agenda", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, response.Code, http.StatusMethodNotAllowed)
	})
	t.Run("calls logger.Info in /agenda http Requests", func(t *testing.T) {
		log.CalledWith = []interface{}{}
		endpoint := "/agenda"
		method := http.MethodPost
		request, _ := http.NewRequest(method, endpoint, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertInsideSlice(t, log.CalledWith, endpoint)
		assertInsideSlice(t, log.CalledWith, method)
	})
}

func TestSessionEndpoint(t *testing.T) {
	server := NewHTTPServer(&log, &aH, &sH)
	t.Run("calls sessionHandler.Post in a /agenda/id/session http POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/agenda/id/session", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertInsideSlice(t, sH.P.CalledWith, response)
		assertInsideSlice(t, sH.P.CalledWith, request)
	})
	t.Run("calls agendaHandler.Get in a /agenda/id/session/id http GET", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/agenda/anID/session/otherID", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertInsideSlice(t, sH.G.CalledWith, response)
		assertInsideSlice(t, sH.G.CalledWith, request)
	})
	t.Run("returns method not allowed for any other method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/agenda/anID/session", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, response.Code, http.StatusMethodNotAllowed)
	})
	t.Run("calls logger.Info in /agenda/id/session/id http Requests", func(t *testing.T) {
		log.CalledWith = []interface{}{}
		endpoint := "/agenda/id/session/id"
		method := http.MethodPost
		request, _ := http.NewRequest(method, endpoint, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertInsideSlice(t, log.CalledWith, endpoint)
		assertInsideSlice(t, log.CalledWith, method)
	})
}

func assertValue(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("want %d, got %d", want, got)
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
