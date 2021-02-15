package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	mathrand "math/rand"
	"time"
)

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"

// GetRandSecretBytes generate random string from a set of alphabet. This function uses CSPRNG
// (crypto-grade random generator) to generate a random value.
func GetRandSecretBytes(n int, alphabet ...byte) ([]byte, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return buf, fmt.Errorf("can't read from buffer: %v", err)
	}

	for i, b := range buf {
		if len(alphabet) == 0 {
			buf[i] = alphanum[b%byte(len(alphanum))]
		} else {
			buf[i] = alphabet[b%byte(len(alphabet))]
		}
	}

	return buf, nil
}

func GetRandSecret(n int, alphabet ...byte) (string, error) {
	buf, err := GetRandSecretBytes(n, alphabet...)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func GetRandBytes(n int, alphabet ...byte) []byte {
	mathrand.Seed(time.Now().UnixNano())
	buf := make([]byte, n)

	for i := range buf {
		if len(alphabet) == 0 {
			//nolint:gosec // Using it for non-sensitive data.
			buf[i] = alphanum[mathrand.Intn(len(alphanum))]
		} else {
			//nolint:gosec // Using it for non-sensitive data.
			buf[i] = alphabet[mathrand.Intn(len(alphabet))]
		}
	}

	return buf
}

func GetRandString(n int, alphabet ...byte) string {
	return string(GetRandBytes(n, alphabet...))
}

func GOBEncodedBytes(val interface{}) (*bytes.Buffer, error) {
	// Use native gob encoding for the fastest serialization.
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return nil, fmt.Errorf("can't gob encode: %v", err)
	}

	return &buf, nil
}

func GOBDecodedBytes(data []byte, val interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(val); err != nil {
		return fmt.Errorf("can't gob decode: %v", err)
	}

	return nil
}
