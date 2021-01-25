package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MaskDSN(t *testing.T) {
	dsns := []struct {
		uri       string
		maskedURI string
		pass      bool
	}{
		{
			uri:       "mongodb://root:password@mongo:27017/db",
			maskedURI: "mongodb://*****:*****@mongo:27017/db",
			pass:      true,
		},
		{
			uri:       "mongodb://root:password@mongo",
			maskedURI: "mongodb://*****:*****@mongo",
			pass:      true,
		},
		{
			uri:  "mongodb://mongo:27017/db",
			pass: false,
		},
		{
			uri:  "",
			pass: false,
		},
	}

	for _, dsn := range dsns {
		actual := maskDSN(dsn.uri)
		if dsn.pass {
			assert.Equal(t, dsn.maskedURI, actual)
			continue
		}

		// No masking
		assert.Equal(t, dsn.uri, actual)
	}
}
