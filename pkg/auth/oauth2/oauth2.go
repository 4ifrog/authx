package oauth2

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/cybersamx/authx/pkg/auth"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/store"
)

// NewAccessToken returns an access token object embedding the token in JWT format.
func NewAccessToken(uid, secrets string, ttl time.Duration) (*models.AccessToken, error) {
	now := time.Now()
	expireAt := time.Now().Add(ttl)
	atID := uuid.New().String()
	jwtToken, err := auth.NewJWT(atID, uid, secrets, now, expireAt)
	if err != nil {
		return nil, err
	}

	at := models.AccessToken{
		ID:       atID,
		Value:    jwtToken,
		UserID:   uid,
		ExpireAt: expireAt,
	}

	return &at, err
}

// NewRefreshToken returns an access refresh token object embedding the token.
func NewRefreshToken(uid string, ttl time.Duration) *models.RefreshToken {
	id := uuid.New().String()
	expireAt := time.Now().Add(ttl)

	rt := models.RefreshToken{
		ID:       id,
		Value:    id,
		UserID:   uid,
		ExpireAt: expireAt,
	}

	return &rt
}

func NewOAuthToken(uid, secrets string, accessTTL, refreshTTL time.Duration) (*models.AccessToken, *models.RefreshToken, error) {
	// Access token
	at, err := NewAccessToken(uid, secrets, accessTTL)
	if err != nil {
		return nil, nil, err
	}

	// Refresh token
	rt := NewRefreshToken(uid, refreshTTL)

	return at, rt, nil
}

func SaveOAuthToken(ctx context.Context, ds store.DataStore, at *models.AccessToken, rt *models.RefreshToken) error {
	if err := ds.SaveAccessToken(ctx, at); err != nil {
		return err
	}
	if err := ds.SaveRefreshToken(ctx, rt); err != nil {
		return err
	}

	return nil
}
