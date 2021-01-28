package auth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/pbkdf2"

	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/store"
)

const (
	pwdHashLen       = 64
	pwdHashIteration = 8
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid authentication credentials")
)

func hashString(str, salt string) string {
	var hashed []byte
	textData := []byte(str)
	saltData := []byte(salt)
	hashed = pbkdf2.Key(textData, saltData, 1<<pwdHashIteration, pwdHashLen, sha256.New)

	return hex.EncodeToString(hashed)
}

func validateHashedString(hashed, clear, salt string) bool {
	hashedClear := hashString(clear, salt)

	return subtle.ConstantTimeCompare([]byte(hashed), []byte(hashedClear)) == 1
}

func Authenticate(parent context.Context, ds store.DataStore, username, password string) (*models.User, error) {
	user, err := ds.GetUserByUsername(parent, username)
	if err == store.ErrorNotFound {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	if user.Username != username || !validateHashedString(user.Password, password, user.Salt) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
