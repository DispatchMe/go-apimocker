package apimocker

import (
	"encoding/json"
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

type Request struct {
	t   *testing.T
	req *http.Request
	w   http.ResponseWriter
}

func NewRequest(t *testing.T, req *http.Request, w http.ResponseWriter) *Request {
	return &Request{t, req, w}
}

func (r *Request) ExpectHeader(k, v string) *Request {
	val := r.req.Header.Get(k)
	if val != v {
		r.t.Errorf("Expected header %s to equal %s (got %s)", k, v, val)
	}
	return r
}

func (r *Request) checkContentType(expected string) bool {
	ct := r.req.Header.Get("Content-Type")
	return ct == expected
}

func (r *Request) Respond() *Response {
	return NewResponse(r.t, r.w)
}

func (r *Request) JSONBody() *httpexpect.Value {
	reporter := httpexpect.NewAssertReporter(r.t)

	if !r.checkContentType("application/json") {
		return nil
	}

	var value interface{}

	if err := json.NewDecoder(r.req.Body).Decode(&value); err != nil {
		r.t.Errorf("Failed to unmarshal JSON body: %s", err.Error())
		return nil
	}
	defer r.req.Body.Close()

	return httpexpect.NewValue(reporter, value)
}
