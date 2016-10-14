package apimocker

import (
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidEndpoint(t *testing.T) {
	server := NewServer(t)
	server.On("GET", "/test", func(req *Request) {
		req.ExpectHeader("Foo", "Bar").Respond().WithStatus(200).Send()
	}).ExpectRequests(1)

	url := server.Start()
	defer server.Stop()

	e := httpexpect.New(t, url)

	e.GET("/test").WithHeader("Foo", "Bar").Expect().Status(200)
	server.AssertExpectations()
}

func TestInvalidEndpoint(t *testing.T) {
	fakeT := new(testing.T)
	server := NewServer(fakeT)
	server.On("GET", "/test", func(req *Request) {
		req.ExpectHeader("Foo", "Baz").Respond().WithStatus(200).Send()
	}).ExpectRequests(1)

	url := server.Start()
	defer server.Stop()

	e := httpexpect.New(t, url)

	e.GET("/test").WithHeader("Foo", "Bar").Expect().Status(200)
	server.AssertExpectations()

	require.Equal(t, true, fakeT.Failed())
}

func TestValidEndpointThatDoesntGetHit(t *testing.T) {
	fakeT := new(testing.T)
	server := NewServer(fakeT)
	server.On("GET", "/test", func(req *Request) {
		req.ExpectHeader("Foo", "Baz").Respond().WithStatus(200).Send()
	}).ExpectRequests(1)

	server.AssertExpectations()

	require.Equal(t, true, fakeT.Failed())
}
