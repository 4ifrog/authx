package storage

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/pbkdf2"
)

const (
	pwdHashLen       = 64
	pwdHashIteration = 8
	alphanum         = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
)

func GetRandString(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	for i, b := range buf {
		buf[i] = alphanum[b%byte(len(alphanum))]
	}

	return string(buf), nil
}

func HashString(str, salt string) string {
	var hashed []byte
	textData := []byte(str)
	saltData := []byte(salt)
	hashed = pbkdf2.Key(textData, saltData, 1<<pwdHashIteration, pwdHashLen, sha256.New)

	return hex.EncodeToString(hashed)
}
