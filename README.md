This package contains helpers for quickly spinning up a fake API server (using `net/http/httptest`), testing the incoming request (using the great JSON tools from `github.com/gavv/httpexpect`), and responding with certain values.

It is meant to be used in tests to mock the "other end" of an HTTP request to get as close to full integration tests as possible without actually relying on an external service.

# Usage
Check the documentation of [https://godoc.org/github.com/gavv/httpexpect#Value](httpexpect.Value) for what to do with the result of `JSONBody()`.

```go

import (
  "testing"
  "github.com/DispatchMe/go-apimocker"
)

func TestMyPackage(t *testing.T) {
  mockAPIServer := apimocker.NewServer(t)

  mockAPIServer.On("GET", "/foo/bar", func(req *apimocker.Request) {
    // Check the header
    req.ExpectHeader("Content-Type", "application/json")

    // Check the body?
    req.JSONBody().Object().Value("some_key").Equal("some_value")

    // Respond
    req.Respond().WithStatus(200).WithHeader("Content-Type", "application/json").WithJSON(apimocker.H{
      "status": "ok",
    }).Send()
  })

  url := mockAPIServer.Start()
  defer mockAPIServer.Stop()

  // Now hit the API!
}
```
