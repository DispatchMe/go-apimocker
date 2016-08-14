package apimocker

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Server struct {
	t          *testing.T
	router     *mux.Router
	testServer *httptest.Server
}

func NewServer(t *testing.T) *Server {
	return &Server{
		t:      t,
		router: mux.NewRouter(),
	}
}

type Handler func(req *Request)

func (s *Server) newHandler(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler(NewRequest(s.t, req, w))
	}
}

func (s *Server) On(method, path string, handler Handler) {
	s.router.HandleFunc(path, s.newHandler(handler)).Methods(method)
}

func (s *Server) Start() string {
	s.testServer = httptest.NewServer(s.router)
	return s.testServer.URL

}

func (s *Server) Stop() {
	s.testServer.Close()
}
