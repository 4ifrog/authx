package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MaskDSN(t *testing.T) {
	tcases := []struct {
		uri      string
		expected string
		pass     bool
	}{
		{
			uri:      "mongodb://root:password@mongo:27017/db",
			expected: "mongodb://*****:*****@mongo:27017/db",
			pass:     true,
		},
		{
			uri:      "mongodb://root:password@mongo",
			expected: "mongodb://*****:*****@mongo",
			pass:     true,
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

	for _, tcase := range tcases {
		actual := maskDSN(tcase.uri)
		if tcase.pass {
			assert.Equal(t, tcase.expected, actual)
			continue
		}

		// No masking
		assert.Equal(t, tcase.uri, actual)
	}
}
