package api

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsAcceptingHTML(t *testing.T) {
	var tcases = []struct {
		accept string
		expect bool
	}{
		{"text/html", true},
		{"application/xhtml+xml", true},
		{"application/xml", true},
		{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif", true},
		{"", false},
		{"q=0.9,image/avif", false},
	}

	for _, tcase := range tcases {
		r := httptest.NewRequest("HEAD", "/", nil)
		r.Header["Accept"] = []string{tcase.accept}

		ok := isAcceptingHTML(r)
		if tcase.expect {
			assert.True(t, ok, `failed due to "%s"`, tcase.accept)
		} else {
			assert.False(t, ok, `failed due to "%s"`, tcase.accept)
		}
	}
}
