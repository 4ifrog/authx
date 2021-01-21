package api

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/storage"
)

const (
	pwdHashLen       = 64
	pwdHashIteration = 8
)

func hashString(str, salt string) string {
	var hashed []byte
	textData := []byte(str)
	saltData := []byte(salt)
	hashed = pbkdf2.Key(textData, saltData, 1<<pwdHashIteration, pwdHashLen, sha256.New)

	return hex.EncodeToString(hashed)
}

func createAccessToken(uid string, ttl int, secrets string) (*models.AccessToken, error) {
	expireAt := time.Now().Add(time.Duration(ttl) * time.Second)
	atID := uuid.New().String()
	atClaims := jwt.MapClaims{}
	atClaims["id"] = atID
	atClaims["user_id"] = uid
	atClaims["expire_at_unix"] = expireAt.Unix()

	atJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	atValue, err := atJWT.SignedString([]byte(secrets))
	if err != nil {
		return nil, err
	}

	at := models.AccessToken{
		ID:       atID,
		Value:    atValue,
		ExpireAt: expireAt,
	}

	return &at, err
}

func createRefreshToken(ttl int) *models.RefreshToken {
	id := uuid.New().String()
	expireAt := time.Now().Add(time.Duration(ttl) * time.Second)

	rt := models.RefreshToken{
		ID:       id,
		Value:    id,
		ExpireAt: expireAt,
	}

	return &rt
}

func createOAuthToken(uid string, cfg *config.Config) (*models.AccessToken, *models.RefreshToken, error) {
	// Access token
	at, err := createAccessToken(uid, cfg.AccessTTL, cfg.AccessSecret)
	if err != nil {
		return nil, nil, err
	}

	// Refresh token
	rt := createRefreshToken(cfg.RefreshTTL)

	return at, rt, nil
}

func saveOAuthToken(ctx context.Context, store storage.Storage, at *models.AccessToken, rt *models.RefreshToken) error {
	if err := store.SaveAccessToken(ctx, at); err != nil {
		return err
	}
	if err := store.SaveRefreshToken(ctx, rt); err != nil {
		return err
	}

	return nil
}

func ValidateHashedString(hashed, clear, salt string) bool {
	hashedClear := hashString(clear, salt)

	return subtle.ConstantTimeCompare([]byte(hashed), []byte(hashedClear)) == 1
}
