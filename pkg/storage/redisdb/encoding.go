package redisdb

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"

	"golang.org/x/crypto/pbkdf2"
)

func gobEncodedBytes(val interface{}) (*bytes.Buffer, error) {
	// Use native gob encoding for the fastest serialization.
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return nil, err
	}

	return &buf, nil
}

func gobDecodedBytes(data []byte, val interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(val); err != nil {
		return err
	}

	return nil
}

func getRandString(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	for i, b := range buf {
		buf[i] = alphanum[b%byte(len(alphanum))]
	}

	return string(buf), nil
}

func hashString(str, salt string) string {
	var hashed []byte
	textData := []byte(str)
	saltData := []byte(salt)
	hashed = pbkdf2.Key(textData, saltData, 1<<pwdHashIteration, pwdHashLen, sha256.New)

	return hex.EncodeToString(hashed)
}

