package apimocker

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type Server struct {
	t          *testing.T
	router     *mux.Router
	testServer *httptest.Server
	endpoints  []*Endpoint
}

func NewServer(t *testing.T) *Server {
	return &Server{
		t:         t,
		router:    mux.NewRouter(),
		endpoints: make([]*Endpoint, 0),
	}
}

type Handler func(req *Request)

func (s *Server) newHandler(handler Handler, endpoint *Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		endpoint.Lock()
		endpoint.hits = endpoint.hits + 1
		endpoint.Unlock()
		handler(NewRequest(s.t, req, w))
	}
}

func (s *Server) On(method, path string, handler Handler) *Endpoint {
	endpoint := new(Endpoint)
	endpoint.Mutex = new(sync.Mutex)
	endpoint.method = method
	endpoint.path = path
	s.endpoints = append(s.endpoints, endpoint)
	s.router.HandleFunc(path, s.newHandler(handler, endpoint)).Methods(method)
	return endpoint
}

func (s *Server) Start() string {
	s.testServer = httptest.NewServer(s.router)
	return s.testServer.URL
}

func (s *Server) Stop() {
	s.testServer.Close()
}

func (s *Server) AssertExpectations() {
	for _, e := range s.endpoints {
		assert.Equal(s.t, e.expectedHits, e.hits, fmt.Sprintf("%s: %s expected %d requests, got %d", e.method, e.path, e.expectedHits, e.hits))
	}
}

type Endpoint struct {
	*sync.Mutex
	path         string
	method       string
	hits         int
	expectedHits int
}

func (e *Endpoint) ExpectRequests(count int) {
	e.Lock()
	e.expectedHits = count
	e.Unlock()
}
