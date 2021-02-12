package api

import (
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type SessionStore struct {
	store sessions.Store
}

type UserSession struct {
	OAuth2Token oauth2.Token
}

const (
	cookieName     = "session"
	cookiePath     = "/"
	cookieMaxAge   = 14 * 24 * 60 * 60
	cookieSecure   = false
	cookieHTTPOnly = false
	keyUserSession = "payload"
)

var (
	ErrSessionSerialization = errors.New("serialization issue with session")
)

func removeSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   cookieName,
		MaxAge: -1, // Remove cookie now.
	}
	http.SetCookie(w, cookie)
}

//nolint:gochecknoinits
func init() {
	// Register the types to serialize the values to the session store.
	gob.Register(new(UserSession))
}

func NewSessionStore(key string) *SessionStore {
	store := sessions.NewCookieStore([]byte(key))
	store.Options = &sessions.Options{
		Path:     cookiePath,
		MaxAge:   cookieMaxAge,
		Secure:   cookieSecure,
		HttpOnly: cookieHTTPOnly,
		SameSite: 0,
	}

	return &SessionStore{
		store: store,
	}
}

func (ss *SessionStore) SetSession(w http.ResponseWriter, r *http.Request, us *UserSession) error {
	session, err := ss.store.Get(r, cookieName)
	if err != nil {
		// If we have trouble getting the cookie, then remove the cookie.
		removeSessionCookie(w)
	}

	session.Values[keyUserSession] = us

	return session.Save(r, w)
}

func (ss *SessionStore) GetSession(r *http.Request) (*UserSession, error) {
	session, err := ss.store.Get(r, cookieName)
	if err != nil {
		return nil, err
	}

	obj := session.Values[keyUserSession]
	if obj == nil {
		return nil, nil
	}
	us, ok := obj.(*UserSession)
	if !ok {
		return nil, ErrSessionSerialization
	}

	return us, nil
}

func (ss *SessionStore) ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, err := ss.store.Get(r, cookieName)
	if err != nil {
		return err
	}

	delete(session.Values, keyUserSession)

	return session.Save(r, w)
}
