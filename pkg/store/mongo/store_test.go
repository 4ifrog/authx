package mongo

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/models"
)

var (
	ds *Store

	longExpiry = time.Now().AddDate(10, 0, 0)
	testUser   = models.User{
		ID:       "user1",
		Username: "test",
		Password: "Password",
		Salt:     "Random",
	}
	testAT = models.AccessToken{
		ID:       "at1",
		Value:    "Random",
		UserID:   "user1",
		ExpireAt: longExpiry,
	}
	testRT = models.RefreshToken{
		ID:       "rt1",
		Value:    "Random",
		UserID:   "user1",
		ExpireAt: longExpiry,
	}
)

func newTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), atomicTimeout)
}

func clearMongo(t *testing.T, s *Store) {
	ctx, cancel := newTestContext()
	defer cancel()

	collections := []string{userCollection, atCollection, rtCollection}
	for _, collect := range collections {
		require.NoError(t, s.db.Collection(collect).Drop(ctx))
	}
}

func seedTestData(t *testing.T, s *Store) {
	ctx, cancel := newTestContext()
	defer cancel()

	values := []struct {
		collection string
		object     interface{}
	}{
		{userCollection, testUser},
		{atCollection, testAT},
		{rtCollection, testRT},
	}

	for _, val := range values {
		_, err := s.db.Collection(val.collection).InsertOne(ctx, val.object)
		require.NoError(t, err)
	}
}

func TestMain(m *testing.M) {
	fmt.Println("Mongo staring up...")
	var code int

	cfg := config.New()
	cfg.MongoAddr = "mongodb://nobody:secrets@localhost:27017/authx"
	if dsn := os.Getenv("AX_MONGO_ADDR"); dsn != "" {
		cfg.MongoAddr = dsn
	}

	ds = New(cfg)
	defer func() {
		fmt.Println("Mongo tearing down...")
		ds.Close()
		os.Exit(code)
	}()

	code = m.Run()
}

func Test_GetUserByID(t *testing.T) {
	clearMongo(t, ds)
	seedTestData(t, ds)

	user, err := ds.GetUser(context.Background(), testUser.ID)

	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, user.ID)
	assert.Equal(t, testUser.Username, user.Username)
}

func Test_GetUserByName(t *testing.T) {
	clearMongo(t, ds)
	seedTestData(t, ds)

	user, err := ds.GetUserByUsername(context.Background(), testUser.Username)

	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, user.ID)
	assert.Equal(t, testUser.Username, user.Username)
}

func Test_SaveUser_UserNotExists(t *testing.T) {
	clearMongo(t, ds)

	ctx := context.Background()
	err := ds.SaveUser(ctx, &testUser)

	assert.NoError(t, err)
	user, err := ds.GetUser(ctx, testUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, user.ID)
	assert.Equal(t, testUser.Username, user.Username)
}

func Test_SaveUser_UserExists(t *testing.T) {
	clearMongo(t, ds)
	seedTestData(t, ds)

	ctx := context.Background()
	err := ds.SaveUser(ctx, &testUser)

	assert.Error(t, err)
}

func Test_RemoveUser_UserExists(t *testing.T) {
	clearMongo(t, ds)
	seedTestData(t, ds)

	ctx := context.Background()
	err := ds.RemoveUser(ctx, testUser.ID)

	assert.NoError(t, err)
	user, err := ds.GetUser(ctx, testUser.ID)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func Test_RemoveUser_UserNotExists(t *testing.T) {
	clearMongo(t, ds)

	ctx := context.Background()
	err := ds.RemoveUser(ctx, testUser.ID)

	assert.NoError(t, err)
}

func Test_SaveAccessToken_TokenNotExists(t *testing.T) {
	clearMongo(t, ds)

	ctx := context.Background()
	err := ds.SaveAccessToken(ctx, &testAT)

	assert.NoError(t, err)
	at, err := ds.GetAccessToken(ctx, testAT.ID)
	assert.NoError(t, err)
	assert.Equal(t, testAT.ID, at.ID)
	assert.Equal(t, testAT.Value, at.Value)
}

func Test_SaveAccessToken_TokenExists(t *testing.T) {
	clearMongo(t, ds)
	seedTestData(t, ds)

	ctx := context.Background()
	err := ds.SaveAccessToken(ctx, &testAT)

	assert.Error(t, err)
}

func Test_RemoveAccessToken_TokenExists(t *testing.T) {
	clearMongo(t, ds)
	seedTestData(t, ds)

	ctx := context.Background()
	err := ds.RemoveAccessToken(ctx, testAT.ID)

	assert.NoError(t, err)
	rt, err := ds.GetRefreshToken(ctx, testRT.ID)
	assert.NoError(t, err)
	assert.Equal(t, testRT.ID, rt.ID)
	assert.Equal(t, testRT.Value, rt.Value)
}

func Test_RemoveAccessToken_TokenNotExists(t *testing.T) {
	clearMongo(t, ds)

	ctx := context.Background()
	err := ds.RemoveAccessToken(ctx, testAT.ID)

	assert.NoError(t, err)
}
