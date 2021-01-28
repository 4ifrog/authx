package redis

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/crypto"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/store"
	"github.com/cybersamx/authx/pkg/utils"
)

const (
	userPrefix         = "users"
	pwdSaltLen         = 24
	accessTokenPrefix  = "access_tokens"
	refreshTokenPrefix = "refresh_tokens"
	redisTimeout       = 15 * time.Second
)

var seedUsers = []struct {
	id       string
	username string
	clearPwd string
}{
	{"0", "admin", "secret"},
	{"1", "chan", "mypassword"},
	{"2", "john", "12345678"},
	{"3", "patel", "patel_rules"},
}

// --- Private Helpers ---

func keyForID(id string) string {
	return fmt.Sprintf("%s:%s", userPrefix, id)
}

func keyForUsername(username string) string {
	return fmt.Sprintf("%s:%s", userPrefix, username)
}

// --- StoreRedis ---

type StoreRedis struct {
	client *redis.Client
}

func New(cfg *config.Config) *StoreRedis {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   0,
	})

	return &StoreRedis{
		client: client,
	}
}

// --- Implements store.DataStore ---

func (sr *StoreRedis) Close() {
	if err := sr.client.Close(); err != nil {
		log.Fatal(err)
	}
}

func (sr *StoreRedis) saveToken(parent context.Context, key string, expireIn time.Duration, buffer *bytes.Buffer) error {
	ctx, cancel := context.WithTimeout(parent, redisTimeout)
	defer cancel()

	if err := sr.client.Set(ctx, key, buffer.Bytes(), expireIn).Err(); err != nil {
		return err
	}

	return nil
}

func (sr *StoreRedis) SaveAccessToken(parent context.Context, at *models.AccessToken) error {
	key := fmt.Sprintf("%s:%s", accessTokenPrefix, at.ID)
	expireIn := time.Until(at.ExpireAt)
	buffer, err := utils.GOBEncodedBytes(at)
	if err != nil {
		return err
	}
	err = sr.saveToken(parent, key, expireIn, buffer)
	if err != nil {
		return err
	}

	return nil
}

func (sr *StoreRedis) SaveRefreshToken(parent context.Context, rt *models.RefreshToken) error {
	// Refresh token
	key := fmt.Sprintf("%s:%s", refreshTokenPrefix, rt.ID)
	expireIn := time.Until(rt.ExpireAt)
	buffer, err := utils.GOBEncodedBytes(rt)
	if err != nil {
		return err
	}
	if err := sr.saveToken(parent, key, expireIn, buffer); err != nil {
		return err
	}

	return nil
}

func (sr *StoreRedis) getUser(parent context.Context, key string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(parent, redisTimeout)
	defer cancel()

	buf, err := sr.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, store.ErrorNotFound
	} else if err != nil {
		return nil, err
	}

	var user models.User
	err = utils.GOBDecodedBytes(buf, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (sr *StoreRedis) SeedUserData() error {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	for _, seedUser := range seedUsers {
		// Generate a user.
		salt, err := utils.GetRandSecret(pwdSaltLen)
		if err != nil {
			panic(err)
		}
		password := crypto.HashString(seedUser.clearPwd, salt)
		user := models.User{
			ID:       seedUser.id,
			Username: seedUser.username,
			Password: password,
			Salt:     salt,
		}

		buffer, err := utils.GOBEncodedBytes(user)
		if err != nil {
			return err
		}

		// Save with ID as the key.
		key := keyForID(user.ID)
		if err := sr.client.Set(ctx, key, buffer.Bytes(), 0).Err(); err != nil {
			return err
		}

		// Save with name as the key.
		key = keyForUsername(user.Username)
		if err := sr.client.Set(ctx, key, buffer.Bytes(), 0).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (sr *StoreRedis) GetUser(parent context.Context, id string) (*models.User, error) {
	return sr.getUser(parent, keyForID(id))
}

func (sr *StoreRedis) GetUserByUsername(parent context.Context, username string) (*models.User, error) {
	return sr.getUser(parent, keyForUsername(username))
}
