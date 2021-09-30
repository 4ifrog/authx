package mongo

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cybersamx/authx/pkg/config"
)

var mongoStore *StoreMongo

func TestMain(m *testing.M) {
	cfg := config.New()
	cfg.MongoAddr = "mongodb://nobody:secrets@localhost:27017/authx"
	if dsn := os.Getenv("AX_MONGO_ADDR"); dsn != "" {
		cfg.MongoAddr = dsn
	}

	var code int
	mongoStore = New(cfg)
	defer func() {
		mongoStore.Close()
		os.Exit(code)
	}()

	code = m.Run()
}

func Test_NewClient(t *testing.T) {
	assert.NotNil(t, mongoStore)
}

func Test_GetUserByID(t *testing.T) {
	user, err := mongoStore.GetUser(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "chan", user.Username)
}

func Test_GetUserByName(t *testing.T) {
	user, err := mongoStore.GetUserByUsername(context.Background(), "chan")

	assert.NoError(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "chan", user.Username)
}
