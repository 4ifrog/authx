package storage

import (
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
	SaveAccessToken(at *models.AccessToken) error
	SaveRefreshToken(rt *models.RefreshToken) error

	// User
	GetUser(id string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	SeedUserData() error
}
