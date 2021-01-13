package api

import (
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

func createTokens(uid string, cfg *config.Config) (*models.AccessToken, *models.RefreshToken, error) {
	// Access token
	expireAt := time.Now().Add(time.Duration(cfg.AccessTTL) * time.Second)
	atID := uuid.New().String()
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = atID
	atClaims["user_id"] = uid
	atClaims["expire_at_unix"] = expireAt.Unix()

	atJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	atValue, err := atJWT.SignedString([]byte(cfg.AccessSecret))
	if err != nil {
		return nil, nil, err
	}

	accessToken := models.AccessToken{
		ID:       atID,
		Value:    atValue,
		ExpireAt: expireAt,
	}

	// Refresh token
	expireAt = time.Now().Add(time.Duration(cfg.RefreshTTL) * time.Second)
	rtID := uuid.New().String()
	rtClaims := jwt.MapClaims{}
	rtClaims["authorized"] = true
	rtClaims["id"] = rtID
	rtClaims["user_id"] = uid
	rtClaims["expire_at_unix"] = expireAt.Unix()

	rtJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rtValue, err := rtJWT.SignedString([]byte(cfg.RefreshSecret))
	if err != nil {
		return nil, nil, err
	}

	refreshToken := models.RefreshToken{
		ID:       uuid.New().String(),
		Value:    rtValue,
		ExpireAt: expireAt,
	}

	return &accessToken, &refreshToken, nil
}

func saveAccessRefreshTokens(store storage.Storage, at *models.AccessToken, rt *models.RefreshToken) error {
	// Access token
	if err := store.SaveAccessToken(at); err != nil {
		return err
	}

	// Refresh token
	if err := store.SaveRefreshToken(rt); err != nil {
		return err
	}

	return nil
}

func ValidateHashedString(hashed, clear, salt string) bool {
	hashedClear := hashString(clear, salt)

	return subtle.ConstantTimeCompare([]byte(hashed), []byte(hashedClear)) == 1
}
