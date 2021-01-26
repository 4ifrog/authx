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
	decoded, err := base64.StdEncoding.DecodeString(chunks[index])
	require.NoError(t, err)
	err = json.Unmarshal(decoded, &payload)
	require.NoError(t, err)
	return httpexpect.NewObject(t, payload)
}

func testExpiry(t *testing.T, at *httpexpect.Object) bool {
	// The expiry in the claims must match the expire_at field in the data payload.

	// expire_at
	expireAt, err := time.Parse(time.RFC3339Nano, at.Value("expire_at").String().Raw())
	require.NoError(t, err)

	// value.exp
	jwtToken := at.Value("value").String().Raw()
	claims := parseChunkFromJWT(t, jwtToken, 1)
	epoch := claims.Value("exp").Number().Raw()

	return expireAt.Unix() == int64(epoch)
}

func testID(t *testing.T, at *httpexpect.Object) bool {
	// The id in the claims must match the user_id field in the data payload.

	// user_id
	userID := at.Value("user_id").String().Raw()

	// value.id
	jwtToken := at.Value("value").String().Raw()
	claims := parseChunkFromJWT(t, jwtToken, 1)
	id := claims.Value("id").String().Raw()

	return userID == id
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
	at := obj.Value("access_token").Object()
	assert.NotEmpty(t, at)
	assert.True(t, testExpiry(t, at))
	assert.True(t, testID(t, at))

	// Validate refresh token
	rt := obj.Value("refresh_token").Object()
	assert.NotNil(t, rt)
	assert.NotEmpty(t, rt.Value("value").String().Raw())
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
