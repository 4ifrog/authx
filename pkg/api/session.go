package api

import (
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type CookieStore struct {
	store sessions.Store
}

type SessionToken struct {
	oauth2.Token
	UserID string
}

// TODO: For some reason setting cookieSecure true is causing headless e2e tests in Docker container to fail. Add TLS?

const (
	cookieName      = "session"
	cookiePath      = "/"
	cookieMaxAge    = 14 * 24 * 60 * 60 // 14 days
	cookieSecure    = false
	cookieHTTPOnly  = true
	keySessionToken = "session-token"
)

var (
	ErrSessionTokenCast     = errors.New("can't cast value as session token type")
	ErrSessionTokenNotFound = errors.New("can't find the session token in the cookie")
)

//nolint:gochecknoinits
func init() {
	// If we set a value of complex type to the session store, gorilla/sessions package will
	// use encoding/gob to serialize/deserialize the value.

	// Register the types to serialize the values to the session store.
	gob.Register(new(SessionToken))
}

func NewCookieStore(key string) *CookieStore {
	store := sessions.NewCookieStore([]byte(key))
	store.Options = &sessions.Options{
		Path:     cookiePath,
		MaxAge:   cookieMaxAge,
		Secure:   cookieSecure,
		HttpOnly: cookieHTTPOnly,
	}

	return &CookieStore{
		store: store,
	}
}

func (cs *CookieStore) SetSessionToken(w http.ResponseWriter, r *http.Request, ut *SessionToken) error {
	// store.Get will always return a session (cookie) even if it's not found.
	session, err := cs.store.Get(r, cookieName)
	if err != nil {
		return err
	}

	// Save token.
	session.Values[keySessionToken] = ut

	// Securely serialized and encoded into a string and then save it to the session (cookie).
	return session.Save(r, w)
}

func (cs *CookieStore) GetSessionToken(r *http.Request) (*SessionToken, error) {
	// store.Get will always return a session (cookie) even if it's not found.
	session, err := cs.store.Get(r, cookieName)
	if err != nil {
		return nil, err
	}

	val := session.Values[keySessionToken]
	if val == nil {
		return nil, ErrSessionTokenNotFound
	}
	ut, ok := val.(*SessionToken)
	if !ok {
		return nil, ErrSessionTokenCast
	}

	return ut, nil
}

func (cs *CookieStore) ClearSessionToken(w http.ResponseWriter, r *http.Request) error {
	// store.Get will always return a session (cookie) even if it's not found.
	session, err := cs.store.Get(r, cookieName)
	if err != nil {
		return err
	}

	delete(session.Values, keySessionToken)

	return session.Save(r, w)
}

func RemoveSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    cookieName,
		Expires: time.Unix(0, 0),
		MaxAge:  -1, // Remove cookie now.
	}

	http.SetCookie(w, cookie)
}
