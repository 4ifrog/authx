package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

const (
	storeKey              = "test-key"
	serializedCookieValue = `MTYzMjU1NjMwNnxEdi1CQkFFQ180SUFBUkFCRUFBQVh2LUNBQUVHYzNSeWFXNW5EQThBRFhObGMzTnBiMjR0ZEc5clpXNFJLbUZ3YVM1VFpYTnphVzl1Vkc5clpXN19nd01CQVF4VFpYTnphVzl1Vkc5clpXNEJfNFFBQVFJQkJWUnZhMlZ1QWYtR0FBRUdWWE5sY2tsRUFRd0FBQUJPXzRVREFRRUZWRzlyWlc0Ql80WUFBUVFCQzBGalkyVnpjMVJ2YTJWdUFRd0FBUWxVYjJ0bGJsUjVjR1VCREFBQkRGSmxabkpsYzJoVWIydGxiZ0VNQUFFR1JYaHdhWEo1QWYtSUFBQUFFUC1IQlFFQkJGUnBiV1VCXzRnQUFBQkJfNFEtQVFFTVlXTmpaWE56TFhSdmEyVnVBUVpDWldGeVpYSUJEWEpsWm5KbGMyZ3RkRzlyWlc0QkR3RUFBQUFPNjY5eUVpZF9FWmotWEFBQkF6RXlNd0E9fIntLK4oBhuXD9Wf01rnjDhzMTEoM4VEKcMTkbC6MuGj` //nolint:lll
)

var (
	longExpiry = time.Now().AddDate(10, 0, 0)
	otoken     = oauth2.Token{
		AccessToken:  "access-token",
		TokenType:    "Bearer",
		RefreshToken: "refresh-token",
		Expiry:       longExpiry,
	}
	us = SessionToken{
		Token:  otoken,
		UserID: "123",
	}
)

func Test_SetSessionToken(t *testing.T) {
	// Request has no cookie.
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	// Create a new cookie and save it in the response.
	store := NewCookieStore(storeKey)
	err := store.SetSessionToken(w, r, &us)
	assert.NoError(t, err)

	defer func() {
		require.NoError(t, w.Result().Body.Close())
	}()
	// Lint rule bodyclose complains even though Close() is there.
	cookies := w.Result().Cookies() //nolint:bodyclose
	assert.Len(t, cookies, 1)
	assert.Equal(t, "session", cookies[0].Name)
	assert.NotEmpty(t, cookies[0].Value)
}

func Test_GetSessionToken_ValidCookie(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	cookie := http.Cookie{
		Name:  "session",
		Value: serializedCookieValue,
	}
	r.AddCookie(&cookie)

	// Get cached cookie.
	store := NewCookieStore(storeKey)
	cachedUS, err := store.GetSessionToken(r)
	assert.NoError(t, err)
	assert.NotNil(t, cachedUS)
	assert.Equal(t, us.UserID, cachedUS.UserID)
	assert.Equal(t, us.Token.AccessToken, cachedUS.Token.AccessToken)
	assert.Equal(t, us.Token.RefreshToken, cachedUS.Token.RefreshToken)
}

func Test_GetSessionToken_InvalidCookie(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	cookie := http.Cookie{
		Name:  "wrong-name",
		Value: "some arbitrary value",
	}
	r.AddCookie(&cookie)

	// Cookie isn't found as it isn't correctly named.
	store := NewCookieStore(storeKey)
	noUS, err := store.GetSessionToken(r)
	assert.Equal(t, ErrSessionTokenNotFound, err)
	assert.Nil(t, noUS)
}

func Test_GetSessionToken_NoCookie(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	// Cookie isn't found as it isn't there.
	store := NewCookieStore(storeKey)
	noUS, err := store.GetSessionToken(r)
	assert.Equal(t, ErrSessionTokenNotFound, err)
	assert.Nil(t, noUS)
}

func Test_ClearSessionToken(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	reqCookie := http.Cookie{
		Name:  "session",
		Value: serializedCookieValue,
	}
	r.AddCookie(&reqCookie)

	// Get cached reqCookie.
	store := NewCookieStore(storeKey)
	cachedUS, err := store.GetSessionToken(r)
	assert.NoError(t, err)
	assert.NotNil(t, cachedUS)

	// Clear the token from the response cookie.
	err = store.ClearSessionToken(w, r)
	assert.NoError(t, err)
	noUS, err := store.GetSessionToken(r)
	assert.Equal(t, ErrSessionTokenNotFound, err)
	assert.Nil(t, noUS)

	defer func() {
		require.NoError(t, w.Result().Body.Close())
	}()
	// Lint rule bodyclose complains even though Close() is there.
	resCookie := w.Result().Cookies()[0] //nolint:bodyclose
	// The request cookie contained the serialized token. After we clear the session, the serialized token
	// should have been cleared, resulting in a smaller payload.
	assert.Less(t, len(resCookie.Value), len(reqCookie.Value))
}
