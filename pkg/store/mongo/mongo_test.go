package mongo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cybersamx/authx/pkg/config"
)

func Test_NewClient(t *testing.T) {
	cfg := config.New()
	cfg.MongoAddr = "mongodb://nobody:secrets@localhost:27017/authx"
	if dsn := os.Getenv("AX_MONGO_ADDR"); dsn != "" {
		cfg.MongoAddr = dsn
	}

	store := New(cfg)
	defer store.Close()

	assert.NotNil(t, store)
}
