package mongo

import (
	"context"
	"log"
	"time"

	"github.com/flowchartsman/retry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/crypto"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/utils"
)

const (
	atomicTimeout  = 15 * time.Second
	pwdSaltLen     = 24
	atCollection   = "access_tokens"
	rtCollection   = "refresh_tokens"
	userCollection = "users"
	database       = "authx"

	// Retry
	retries      = 5
	initialDelay = 1 * time.Second
	maxDelay     = 6 * time.Second
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

// setupMongo configure the mongo store with indexes, collections, etc.
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

// --- StoreMongo ---

type StoreMongo struct {
	client *mongo.Client
	db     *mongo.Database
}

func newClient(parent context.Context, dsn string, retries int, initialDelay, maxDelay time.Duration) (*mongo.Client, error) {
	var client *mongo.Client
	timeout := maxDelay

	retrier := retry.NewRetrier(retries, initialDelay, maxDelay)
	err := retrier.RunContext(parent, func(ctx context.Context) (retErr error) {
		var cterr error
		client, cterr = mongo.NewClient(options.Client().ApplyURI(dsn))
		if cterr != nil {
			log.Fatal("can't create an instance of mongo client")
		}

		// Disconnect only if we can't connect or ping the store.
		closeFn := func() {
			dctx, dcancel := context.WithTimeout(ctx, timeout)
			defer dcancel()
			if derr := client.Disconnect(dctx); derr != nil && derr != mongo.ErrClientDisconnected {
				log.Printf("can't close connection to mongo: %v\n", derr)
				retErr = derr
			}
		}

		// Connect the store.
		log.Printf("attempting to connect mongo %s\n", maskDSN(dsn))
		cctx, ccancel := context.WithTimeout(ctx, timeout)
		defer ccancel()
		if cerr := client.Connect(cctx); cerr != nil {
			defer closeFn()

			if cerr == topology.ErrTopologyConnected {
				// Already connected, so continue to ping the store.
			} else {
				log.Printf("can't connect mongo: %v\n", cerr)
				retErr = cerr
				return
			}
		}

		// Ping the store.
		log.Printf("attempting to ping mongo %s\n", maskDSN(dsn))
		pctx, pcancel := context.WithTimeout(ctx, timeout)
		defer pcancel()
		if perr := client.Ping(pctx, readpref.Primary()); perr != nil {
			defer closeFn()
			log.Printf("can't ping mongo: %v\n", perr)
			retErr = perr
			return
		}

		retErr = nil
		return
	})

	return client, err
}

func New(cfg *config.Config) *StoreMongo {
	dsn := cfg.MongoAddr

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := newClient(ctx, dsn, retries, initialDelay, maxDelay)

	if err != nil {
		log.Printf("exhausted all %d retries", retries)
		panic(err)
	}

	log.Printf("connected to mongo %s successfully\n", maskDSN(dsn))

	// Additional setup.
	db := client.Database(database)
	if err := setupMongo(ctx, db); err != nil {
		panic(err)
	}

	return &StoreMongo{
		client: client,
		db:     db,
	}
}

// --- Implements store.Storage ---

func (sm *StoreMongo) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), atomicTimeout)
	defer cancel()
	if err := sm.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (sm *StoreMongo) SaveAccessToken(parent context.Context, at *models.AccessToken) error {
	ctx, cancel := context.WithTimeout(parent, atomicTimeout)
	defer cancel()

	_, err := sm.db.Collection(atCollection).InsertOne(ctx, at)
	if err != nil {
		return err
	}

	return nil
}

func (sm *StoreMongo) SaveRefreshToken(parent context.Context, rt *models.RefreshToken) error {
	ctx, cancel := context.WithTimeout(parent, atomicTimeout)
	defer cancel()

	_, err := sm.db.Collection(rtCollection).InsertOne(ctx, rt)
	if err != nil {
		return err
	}

	return nil
}

func (sm *StoreMongo) getUser(parent context.Context, key, val string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(parent, atomicTimeout)
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

func (sm *StoreMongo) SeedUserData() error {
	ctx, cancel := context.WithTimeout(context.Background(), atomicTimeout)
	defer cancel()

	if err := sm.db.Collection(userCollection).Drop(ctx); err != nil {
		panic(err)
	}

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

		_, err = sm.db.Collection(userCollection).InsertOne(ctx, user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sm *StoreMongo) GetUser(parent context.Context, id string) (*models.User, error) {
	return sm.getUser(parent, "_id", id)
}

func (sm *StoreMongo) GetUserByUsername(parent context.Context, username string) (*models.User, error) {
	return sm.getUser(parent, "username", username)
}
