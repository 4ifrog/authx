package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/storage"
)

const (
	atomicTimeout  = 30 * time.Second
	pwdSaltLen     = 24
	atCollection   = "access_tokens"
	rtCollection   = "refresh_tokens"
	userCollection = "users"
	database       = "authx"
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

func setupMongo(parentCtx context.Context, db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(parentCtx, atomicTimeout)
	defer cancel()

	for _, colName := range []string{atCollection, rtCollection} {
		// Create TTL index.
		opts := IndexOptions{isTTL: true}
		_, err := createIndex(ctx, db.Collection(colName), "expireAt", &opts)
		if err != nil {
			return err
		}
	}

	return nil
}

// --- StorageMongo ---

type StorageMongo struct {
	client *mongo.Client
	db     *mongo.Database
}

func New(cfg *config.Config) *StorageMongo {
	uri := cfg.MongoAddr

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		panic(err)
	}

	db := client.Database(database)

	// Additional setup on the collections.
	if err := setupMongo(ctx, db); err != nil {
		panic(err)
	}

	return &StorageMongo{
		client: client,
		db:     db,
	}
}

// --- Implements storage.Storage ---

func (sm *StorageMongo) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), atomicTimeout)
	defer cancel()
	if err := sm.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (sm *StorageMongo) SaveAccessToken(at *models.AccessToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), atomicTimeout)
	defer cancel()

	_, err := sm.db.Collection(atCollection).InsertOne(ctx, at)
	if err != nil {
		return err
	}

	return nil
}

func (sm *StorageMongo) SaveRefreshToken(rt *models.RefreshToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), atomicTimeout)
	defer cancel()

	_, err := sm.db.Collection(rtCollection).InsertOne(ctx, rt)
	if err != nil {
		return err
	}

	return nil
}

func (sm *StorageMongo) getUser(key, val string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), atomicTimeout)
	defer cancel()

	var user models.User
	err := sm.db.Collection(userCollection).FindOne(ctx, bson.D{
		{Key: key, Value: val},
	}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (sm *StorageMongo) SeedUserData() error {
	ctx, cancel := context.WithTimeout(context.Background(), atomicTimeout)
	defer cancel()

	for _, seedUser := range seedUsers {
		// Generate a user.
		salt, err := storage.GetRandString(pwdSaltLen)
		if err != nil {
			panic(err)
		}
		password := storage.HashString(seedUser.clearPwd, salt)
		user := models.User{
			ID:       seedUser.id,
			Username: seedUser.username,
			Password: password,
			Salt:     salt,
		}

		_, err = sm.db.Collection(userCollection).InsertOne(ctx, user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sm *StorageMongo) GetUser(id string) (*models.User, error) {
	return sm.getUser("_id", id)
}

func (sm *StorageMongo) GetUserByUsername(username string) (*models.User, error) {
	return sm.getUser("username", username)
}
