package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ErrorPage401(t *testing.T) {
	// Setup
	expect := newHTTPExpect(t)

	// Run
	obj := expect.GET("/v1/userinfo").
		Expect().
		Status(http.StatusUnauthorized)

	assert.NotNil(t, obj)
}
