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
	keySession     = "session"
	keyPayload     = "payload"
	cookiePath     = "/"
	cookieMaxAge   = 14 * 24 * 60 * 60
	cookieSecure   = false
	cookieHTTPOnly = false
)

var (
	ErrSessionSerialization = errors.New("serialization issue with session")
)

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
	session, err := ss.store.Get(r, keySession)
	if err != nil {
		return err
	}

	session.Values[keyPayload] = us

	return session.Save(r, w)
}

func (ss *SessionStore) GetSession(r *http.Request) (*UserSession, error) {
	session, err := ss.store.Get(r, keySession)
	if err != nil {
		return nil, err
	}

	obj := session.Values[keyPayload]
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
	session, err := ss.store.Get(r, keySession)
	if err != nil {
		return err
	}

	delete(session.Values, keyPayload)

	return session.Save(r, w)
}
