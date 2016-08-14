package apimocker

import (
	"encoding/json"
	"net/http"
	"testing"
)

type Response struct {
	t *testing.T

	body []byte
	resp http.ResponseWriter
}

type H map[string]interface{}

func NewResponse(t *testing.T, resp http.ResponseWriter) *Response {
	return &Response{
		t:    t,
		resp: resp,
	}
}

func (r *Response) WithStatus(code int) *Response {
	r.resp.WriteHeader(code)
	return r
}

func (r *Response) WithHeader(k, v string) *Response {
	r.resp.Header().Set(k, v)
	return r
}

func (r *Response) WithJSON(obj interface{}) *Response {
	r.WithHeader("Content-Type", "application/json")

	body, err := json.Marshal(obj)
	if err != nil {
		r.t.Errorf("Failed to marshal JSON response into JSON: %s", err.Error())
	} else {
		r.body = body
	}
	return r
}

func (r *Response) Send() {
	r.resp.Write(r.body)
}
