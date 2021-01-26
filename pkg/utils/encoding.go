package utils

import (
	"crypto/rand"
	mathrand "math/rand"
	"time"
)

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"

// GetRandSecretBytes generate random string from a set of alphabet. This function uses CSPRNG
// (crypto-grade random generator) to generate a random value.
func GetRandSecretBytes(n int, alphabet ...byte) ([]byte, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return buf, err
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
			buf[i] = alphanum[mathrand.Intn(len(alphanum))] //nolint:gosec
		} else {
			buf[i] = alphabet[mathrand.Intn(len(alphabet))] //nolint:gosec
		}
	}

	return buf
}

func GetRandString(n int, alphabet ...byte) string {
	return string(GetRandBytes(n, alphabet...))
}
