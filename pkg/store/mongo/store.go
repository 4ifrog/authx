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
	"github.com/cybersamx/authx/pkg/store"
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
func setupMongo(parent context.Context, db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(parent, atomicTimeout)
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

// --- Store ---

type Store struct {
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

func New(cfg *config.Config) *Store {
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

	return &Store{
		client: client,
		db:     db,
	}
}

func (s *Store) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), atomicTimeout)
	defer cancel()
	if err := s.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (s *Store) SeedUserData() error {
	ctx, cancel := context.WithTimeout(context.Background(), atomicTimeout)
	defer cancel()

	if err := s.db.Collection(userCollection).Drop(ctx); err != nil {
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

		_, err = s.db.Collection(userCollection).InsertOne(ctx, user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) getAndBindObject(parent context.Context, collection, key, val string, obj interface{}) error {
	ctx, cancel := context.WithTimeout(parent, atomicTimeout)
	defer cancel()

	err := s.db.Collection(collection).FindOne(ctx, bson.D{
		{Key: key, Value: val},
	}).Decode(obj)
	if err == mongo.ErrNoDocuments {
		return store.ErrorNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func (s *Store) saveObject(parent context.Context, collection string, obj interface{}) error {
	ctx, cancel := context.WithTimeout(parent, atomicTimeout)
	defer cancel()

	_, err := s.db.Collection(collection).InsertOne(ctx, obj)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) removeObject(parent context.Context, collection, id string) error {
	ctx, cancel := context.WithTimeout(parent, atomicTimeout)
	defer cancel()

	_, err := s.db.Collection(collection).DeleteOne(ctx, bson.D{
		{Key: "_id", Value: id},
	})

	return err
}

func (s *Store) GetUser(parent context.Context, id string) (*models.User, error) {
	var user models.User
	if err := s.getAndBindObject(parent, userCollection, "_id", id, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserByUsername(parent context.Context, username string) (*models.User, error) {
	var user models.User
	if err := s.getAndBindObject(parent, userCollection, "username", username, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) SaveUser(parent context.Context, user *models.User) error {
	return s.saveObject(parent, userCollection, user)
}

func (s *Store) RemoveUser(parent context.Context, id string) error {
	return s.removeObject(parent, userCollection, id)
}

func (s *Store) GetAccessToken(parent context.Context, id string) (*models.AccessToken, error) {
	var at models.AccessToken
	if err := s.getAndBindObject(parent, atCollection, "_id", id, &at); err != nil {
		return nil, err
	}

	return &at, nil
}

func (s *Store) SaveAccessToken(parent context.Context, at *models.AccessToken) error {
	return s.saveObject(parent, atCollection, at)
}

func (s *Store) RemoveAccessToken(parent context.Context, id string) error {
	return s.removeObject(parent, atCollection, id)
}

func (s *Store) GetRefreshToken(parent context.Context, id string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	if err := s.getAndBindObject(parent, rtCollection, "_id", id, &rt); err != nil {
		return nil, err
	}

	return &rt, nil
}

func (s *Store) SaveRefreshToken(parent context.Context, rt *models.RefreshToken) error {
	return s.saveObject(parent, rtCollection, rt)
}

func (s *Store) RemoveRefreshToken(parent context.Context, id string) error {
	return s.removeObject(parent, rtCollection, id)
}
