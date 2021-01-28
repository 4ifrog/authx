package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func parseChunkFromJWT(t *testing.T, jwt string, index int) *httpexpect.Object {
	chunks := strings.Split(jwt, ".")
	require.Equal(t, 3, len(chunks))
	var payload map[string]interface{}
	decoded, err := base64.RawStdEncoding.DecodeString(chunks[index])
	require.NoError(t, err)
	err = json.Unmarshal(decoded, &payload)
	require.NoError(t, err)
	return httpexpect.NewObject(t, payload)
}

func testExpiry(t *testing.T, obj *httpexpect.Object) bool {
	// The expiry in the claims must match the expire_at field in the data payload.

	// expiry
	expireAt, err := time.Parse(time.RFC3339Nano, obj.Value("expiry").String().Raw())
	require.NoError(t, err)

	// exp field in the base64 decoded JWT
	jwtToken := obj.Value("access_token").String().Raw()
	claims := parseChunkFromJWT(t, jwtToken, 1)
	epoch := claims.Value("exp").Number().Raw()

	return expireAt.Unix() == int64(epoch)
}

func Test_PostSignIn(t *testing.T) {
	// Setup
	server := httptest.NewServer(a.Router)
	expect := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  server.URL,
		Reporter: httpexpect.NewRequireReporter(t),
	})

	// Run
	obj := expect.POST("/v1/signin").WithBytes([]byte(`{"username": "chan", "password": "mypassword"}`)).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// Validate access token
	assert.NotNil(t, obj)
	at := obj.Value("access_token").String().Raw()
	assert.NotEmpty(t, at)
	assert.True(t, testExpiry(t, obj))

	// Validate refresh token
	rt := obj.Value("refresh_token").String().Raw()
	assert.NotEmpty(t, rt)
}

func Test_SignOutHandler(t *testing.T) {
	// Setup
	expect := newHTTPExpect(t)

	// Add cookie
	req := expect.GET("/v1/signout")

	// Run
	msg := req.Expect().Status(http.StatusOK).Body().Raw()

	// Validate
	assert.Equal(t, `"logout"`, msg)
}
