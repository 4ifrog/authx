package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cybersamx/authx/pkg/auth"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testUser = models.User{
		ID:       "1",
		Username: "chan",
		Password: "mypassword",
	}
)

func newTestAccessToken(t *testing.T, cfg *config.Config, targetDate time.Time) *models.AccessToken {
	distantTTL := time.Until(targetDate)
	at, err := auth.NewAccessToken(testUser.ID, cfg.AccessSecret, distantTTL)
	assert.NoError(t, err)

	return at
}

func newValidAccessToken(t *testing.T, cfg *config.Config) *models.AccessToken {
	return newTestAccessToken(t, cfg, time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC))
}

func newExpiredAccessToken(t *testing.T, cfg *config.Config) *models.AccessToken {
	return newTestAccessToken(t, cfg, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
}

func newHTTPExpect(t *testing.T) *httpexpect.Expect {
	srv := httptest.NewServer(testapp.Router)

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
	expect := newHTTPExpect(t)

	// Run
	userJSON, err := json.Marshal(testUser)
	require.NoError(t, err)
	obj := expect.POST("/v1/signin").
		WithBytes(userJSON).
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

	// Run
	req := expect.POST("/v1/signout")
	msg := req.Expect().Status(http.StatusOK).Body().Raw()

	// Validate
	assert.Equal(t, `"no access token"`, msg)
}

func Test_ProfileHandler_ValidAccessToken(t *testing.T) {
	// Setup
	expect := newHTTPExpect(t)

	// Run
	at := newValidAccessToken(t, testapp.Config)
	obj := expect.GET("/v1/userinfo").
		WithHeader("Authorization", fmt.Sprintf("Bearer %s", at.Value)).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// Validate
	assert.NotNil(t, obj)
	assert.Equal(t, obj.Value("username").String().Raw(), "chan")
	assert.Equal(t, obj.Value("id").String().Raw(), "1")
}

func Test_ProfileHandler_ExpiredAccessToken(t *testing.T) {
	// Setup
	expect := newHTTPExpect(t)

	// Run
	at := newExpiredAccessToken(t, testapp.Config)
	expect.GET("/v1/userinfo").
		WithHeader("Authorization", fmt.Sprintf("Bearer %s", at.Value)).
		Expect().
		Status(http.StatusUnauthorized)
}
