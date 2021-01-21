package storage

import (
	"context"
	"errors"

	"github.com/cybersamx/authx/pkg/models"
)

var (
	ErrorNotFound = errors.New("user not found")
)

type Storage interface {
	// General
	Close()

	// Token
	SaveAccessToken(parent context.Context, at *models.AccessToken) error
	SaveRefreshToken(parent context.Context, rt *models.RefreshToken) error

	// User
	GetUser(parent context.Context, id string) (*models.User, error)
	GetUserByUsername(parent context.Context, username string) (*models.User, error)
	SeedUserData() error
}
