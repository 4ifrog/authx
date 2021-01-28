package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/store"
)

const (
	tokenType = "Bearer"
)

// NewAccessToken returns an access token object embedding the token in JWT format.
func NewAccessToken(uid, secrets string, ttl time.Duration) (*models.AccessToken, error) {
	now := time.Now()
	expireAt := time.Now().Add(ttl)
	atID := uuid.New().String()
	jwtToken, err := NewJWT(atID, uid, secrets, now, expireAt)
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

type CreateOAuthTokenParams struct {
	UID        string
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func CreateOAuthToken(parent context.Context, ds store.DataStore, uid, secret string, aTTL, rTTL time.Duration) (*oauth2.Token, error) {
	// Access token
	at, err := NewAccessToken(uid, secret, aTTL)
	if err != nil {
		return nil, err
	}

	// Refresh token
	rt := NewRefreshToken(uid, rTTL)

	// Save the tokens
	if err := ds.SaveAccessToken(parent, at); err != nil {
		return nil, err
	}
	if err := ds.SaveRefreshToken(parent, rt); err != nil {
		return nil, err
	}

	// OAuth2 token
	otoken := oauth2.Token{
		AccessToken:  at.Value,
		TokenType:    tokenType,
		RefreshToken: rt.Value,
		Expiry:       at.ExpireAt,
	}

	return &otoken, nil
}
