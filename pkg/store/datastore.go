package store

import (
	"context"
	"errors"

	"github.com/cybersamx/authx/pkg/models"
)

var (
	ErrorNotFound = errors.New("object not found")
)

type DataStore interface {
	Close()

	GetUser(parent context.Context, id string) (*models.User, error)
	GetUserByUsername(parent context.Context, username string) (*models.User, error)
	SaveUser(parent context.Context, user *models.User) error
	RemoveUser(parent context.Context, id string) error

	GetAccessToken(parent context.Context, id string) (*models.AccessToken, error)
	SaveAccessToken(parent context.Context, at *models.AccessToken) error
	RemoveAccessToken(parent context.Context, id string) error

	GetRefreshToken(parent context.Context, id string) (*models.RefreshToken, error)
	SaveRefreshToken(parent context.Context, rt *models.RefreshToken) error
	RemoveRefreshToken(parent context.Context, id string) error
}
