package api

import (
	"time"

	"github.com/cybersamx/authx/pkg/models"
)

type AccessToken struct {
	Value    string    `json:"value"`
	UserID   string    `json:"user_id"`
	ExpireAt time.Time `json:"expire_at"`
}

type RefreshToken struct {
	Value    string    `json:"value"`
	UserID   string    `json:"user_id"`
	ExpireAt time.Time `json:"expire_at"`
}

type OAuth2Tokens struct {
	AccessToken  *AccessToken  `json:"access_token"`
	RefreshToken *RefreshToken `json:"refresh_token"`
}

func NewAccessTokenModel(model *models.AccessToken) *AccessToken {
	return &AccessToken{
		Value:    model.Value,
		UserID:   model.UserID,
		ExpireAt: model.ExpireAt,
	}
}

func NewRefreshTokenModel(model *models.RefreshToken) *RefreshToken {
	return &RefreshToken{
		Value:    model.Value,
		UserID:   model.UserID,
		ExpireAt: model.ExpireAt,
	}
}
