package api

import (
	"context"
	"fmt"
	"time"

	"github.com/cybersamx/authx/pkg/crypto"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/store"
	"github.com/cybersamx/authx/pkg/utils"
)

const (
	storeTimeout = 15 * time.Second
	pwdSaltLen   = 24
)

var (
	seedUsers = []struct {
		id       string
		username string
		clearPwd string
	}{
		{"0", "admin", "secret"},
		{"1", "chan", "mypassword"},
		{"2", "john", "12345678"},
		{"3", "patel", "patel_rules"},
	}
)

func SeedUserData(ds store.DataStore) error {
	ctx, cancel := context.WithTimeout(context.Background(), storeTimeout)
	defer cancel()

	for _, seedUser := range seedUsers {
		_, err := ds.GetUser(ctx, seedUser.id)
		if err == store.ErrorNotFound {
			// Generate a user.
			salt, errr := utils.GetRandSecret(pwdSaltLen)
			if errr != nil {
				panic(err)
			}
			password := crypto.HashString(seedUser.clearPwd, salt)
			user := models.User{
				ID:       seedUser.id,
				Username: seedUser.username,
				Password: password,
				Salt:     salt,
			}

			if errr := ds.SaveUser(ctx, &user); errr != nil {
				return fmt.Errorf("failed to save seed user: %v", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to get user: %v", err)
		}
	}

	return nil
}
