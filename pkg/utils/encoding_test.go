package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func containInString(search string, alphabet []byte) bool {
	for _, s := range []byte(search) {
		for _, c := range alphabet {
			if s == c {
				return true
			}
		}
	}

	return false
}

func containInBytes(search, alphabet []byte) bool {
	for _, s := range search {
		for _, a := range alphabet {
			if s == a {
				return true
			}
		}
	}

	return false
}

func TestRandomSecret(t *testing.T) {
	cases := []struct {
		size     int
		alphabet []byte
	}{
		{size: 1, alphabet: []byte("abcdefgh")},
		{size: 1, alphabet: []byte("12387ASDFCMeurye")},
		{size: 5, alphabet: []byte("d338fkaeruf54")},
		{size: 5, alphabet: []byte("a4")},
		{size: 10, alphabet: []byte("3$328347dwerh")},
		{size: 20, alphabet: []byte("383wer834-9fsfj#34")},
	}

	for _, tcase := range cases {
		val := tcase

		t.Run("GetRandSecretBytes", func(t *testing.T) {
			// Get random string.
			randString, err := GetRandSecretBytes(val.size, val.alphabet...)

			// Validate.
			assert.NoError(t, err)
			assert.Equal(t, val.size, len(randString))
			assert.True(t, containInBytes(randString, val.alphabet))
		})

		t.Run("GetRandSecret", func(t *testing.T) {
			// Get random string.
			randString, err := GetRandSecret(val.size, val.alphabet...)

			// Validate.
			assert.NoError(t, err)
			assert.Equal(t, val.size, len(randString))
			assert.True(t, containInString(randString, val.alphabet))
		})
	}
}

func TestRandomString(t *testing.T) {
	cases := []struct {
		size     int
		alphabet []byte
	}{
		{size: 1, alphabet: []byte("abcdefgh")},
		{size: 1, alphabet: []byte("12387ASDFCMeurye")},
		{size: 5, alphabet: []byte("d338fkaeruf54")},
		{size: 5, alphabet: []byte("a4")},
		{size: 10, alphabet: []byte("3$328347dwerh")},
		{size: 20, alphabet: []byte("383wer834-9fsfj#34")},
	}

	for _, tcase := range cases {
		val := tcase

		t.Run("GetRandSecretBytes", func(t *testing.T) {
			// Get random string.
			randString := GetRandBytes(val.size, val.alphabet...)

			// Validate.
			assert.Equal(t, val.size, len(randString))
			assert.True(t, containInBytes(randString, val.alphabet))
		})

		t.Run("GetRandSecret", func(t *testing.T) {
			// Get random string.
			randString := GetRandString(val.size, val.alphabet...)

			// Validate.
			assert.Equal(t, val.size, len(randString))
			assert.True(t, containInString(randString, val.alphabet))
		})
	}
}
