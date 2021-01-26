package crypto

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/pbkdf2"
)

const (
	pwdHashLen       = 64
	pwdHashIteration = 8
)

func HashString(str, salt string) string {
	var hashed []byte
	textData := []byte(str)
	saltData := []byte(salt)
	hashed = pbkdf2.Key(textData, saltData, 1<<pwdHashIteration, pwdHashLen, sha256.New)

	return hex.EncodeToString(hashed)
}
