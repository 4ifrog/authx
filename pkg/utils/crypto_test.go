package utils

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cases = []struct {
	description string
	clear       string
	md5         string
	sha1        string
}{
	{
		description: "text fewer than 32-character",
		clear:       "A",
		md5:         "7fc56270e7a70fa81a5935b72eacbe29",
		sha1:        "6dcd4ce23d88e2ee9568ba546c007c63d9131c1b",
	},
	{
		description: "text exactly 32-character",
		clear:       "abcdefghijklmnopqrstuvwxyz123456",
		md5:         "9ee779cd2abcde48524485572c6ce2a2",
		sha1:        "18fdf21db6a592dac66e4bcefad19d6b8d0428f4",
	},
	{
		description: "text more than 32-character",
		clear:       "ZDURSsdf_wirevmwr83dcwer2384972#@@vmzepowrpf3",
		md5:         "004778972005af384ca2048e7c250005",
		sha1:        "3b79e1c12804138865fbbf6ff0544d02ded8da0f",
	},
}

//nolint:dupl // There's enough variation that it should be marked as a dupe.
func TestMD5(t *testing.T) {
	for _, tcase := range cases {
		tcase := tcase
		desc := fmt.Sprintf("MD5 - %s", tcase.description)
		t.Run(desc, func(t *testing.T) {
			// Test MD5() - positive
			hash := MD5(tcase.clear)
			assert.Equal(t, tcase.md5, hash)

			// Test MD5() - negative
			randStr, err := GetRandSecret(3)
			require.NoError(t, err)
			tampered := tcase.clear + randStr
			hash = MD5(tampered)
			assert.NotEqual(t, tcase.md5, hash)

			// Test MD5Bytes() - positive
			hashBytes := MD5Bytes(tcase.clear)
			assert.Equal(t, tcase.md5, hex.EncodeToString(hashBytes))

			// Test MD5Bytes() - negative
			randStr, err = GetRandSecret(3)
			require.NoError(t, err)
			tampered = tcase.clear + randStr
			hashBytes = MD5Bytes(tampered)
			assert.NotEqual(t, tcase.md5, hex.EncodeToString(hashBytes))
		})
	}
}

//nolint:dupl // There's enough variation that it should be marked as a dupe.
func TestSHA1(t *testing.T) {
	for _, tcase := range cases {
		tcase := tcase
		desc := fmt.Sprintf("SHA1 - %s", tcase.description)
		t.Run(desc, func(t *testing.T) {
			// Test SHA1() - positive
			hash := SHA1(tcase.clear)
			assert.Equal(t, tcase.sha1, hash)

			// Test SHA1() - negative
			randStr, err := GetRandSecret(3)
			require.NoError(t, err)
			tampered := tcase.clear + randStr
			hash = SHA1(tampered)
			assert.NotEqual(t, tcase.sha1, hash)

			// Test SHA1Bytes() - positive
			hashBytes := SHA1Bytes(tcase.clear)
			assert.Equal(t, tcase.sha1, hex.EncodeToString(hashBytes))

			// Test SHA1Bytes() - negative
			randStr, err = GetRandSecret(3)
			require.NoError(t, err)
			tampered = tcase.clear + randStr
			hashBytes = SHA1Bytes(tampered)
			assert.NotEqual(t, tcase.sha1, hex.EncodeToString(hashBytes))
		})
	}
}
