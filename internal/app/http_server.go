package server

import (
	"net/http"
	"regexp"

	"github.com/cesarFuhr/gocrypto/internal/app/ports"
)

// HTTPServer http server interface
type HTTPServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// HTTPLogger http server logger
type HTTPLogger interface {
	Info(...interface{})
}

type httpServer struct {
	routes []*route
}

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

// NewHTTPServer creates a new http handler
func NewHTTPServer(
	l HTTPLogger,
	aH ports.AgendaHandler,
) HTTPServer {
	logger := newLoggerMiddleware(l)
	routes := []*route{
		createRoute("/agenda/{0,}$", logger(handleCreateAgenda(aH))),
		createRoute("/agenda/[^/]{0,}/{0,}$", logger(handleFindAgenda(aH))),
	}
	return &httpServer{
		routes: routes,
	}
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range s.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func handleCreateAgenda(h ports.AgendaHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Post(w, r)
			return
		}
		methodNotAllowed(w, r)
	})
}

func handleFindAgenda(h ports.AgendaHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.Get(w, r)
			return
		}
		methodNotAllowed(w, r)
	})
}

func createRoute(pattern string, handler http.Handler) *route {
	rg := regexp.MustCompile(pattern)
	return &route{
		pattern: rg,
		handler: handler,
	}
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte{})
}
