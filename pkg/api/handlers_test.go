package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
)

func newHTTPExpect(t *testing.T) *httpexpect.Expect {
	srv := httptest.NewServer(a.Router)

	cfg := httpexpect.Config{
		BaseURL:  srv.URL,
		Reporter: httpexpect.NewRequireReporter(t),
		Printers: nil,
	}

	return httpexpect.WithConfig(cfg)
}

func Test_PostSignIn(t *testing.T) {
	server := httptest.NewServer(a.Router)

	expect := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  server.URL,
		Reporter: httpexpect.NewRequireReporter(t),
	})

	obj := expect.POST("/v1/signin").WithBytes([]byte(`{"username": "chan", "password": "mypassword"}`)).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.NotEmpty(t, obj.Value("access_token").String().Raw())
	assert.NotEmpty(t, obj.Value("refresh_token").String().Raw())
}

func Test_SignOutHandler(t *testing.T) {
	// Setup
	expect := newHTTPExpect(t)

	req := expect.GET("/v1/signout")
	req.WithHeaders(map[string]string{
		"Authorization": "Bearer XXXX",
	})

	// Run
	msg := req.Expect().Status(http.StatusNotImplemented).Body().Raw()

	// Validate
	assert.Equal(t, `"Not implemented"`, msg)
}
