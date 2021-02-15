package utils

import (
	//nolint:gosec // Need MD5 for encoding text for querying a service
	"crypto/md5"
	//nolint:gosec // Need SHA1 for encoding text for querying a service
	"crypto/sha1"
	"encoding/hex"
)

// MD5 encodes a string to a MD5 encoded string.
func MD5(str string) string {
	return hex.EncodeToString(MD5Bytes(str))
}

// MD5Bytes encodes a string to MD5 bytes.
func MD5Bytes(str string) []byte {
	//nolint:gosec
	hash := md5.New()
	_, err := hash.Write([]byte(str))
	if err != nil {
		return []byte{}
	}
	return hash.Sum(nil)
}

// SHA1 encodes a string to a 40-byte long SHA1 encoded string.
func SHA1(str string) string {
	return hex.EncodeToString(SHA1Bytes(str))
}

// SHA1Bytes encodes a string to SHA1 bytes.
func SHA1Bytes(str string) []byte {
	//nolint:gosec
	hash := sha1.New()
	_, err := hash.Write([]byte(str))
	if err != nil {
		return []byte{}
	}
	return hash.Sum(nil)
}
