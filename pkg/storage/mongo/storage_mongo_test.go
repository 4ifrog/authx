package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cybersamx/authx/pkg/config"
)

func TestNewClient(t *testing.T) {
	cfg := config.New()
	cfg.MongoAddr = "mongodb://nobody:secrets@localhost:27017/authx"

	store := New(cfg)
	defer store.Close()

	assert.NotNil(t, store)
}
