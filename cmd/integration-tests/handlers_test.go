package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
)

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
