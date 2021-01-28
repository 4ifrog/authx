package auth

import (
	"errors"
	"time"

	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

// Types of claims

// JWTClaims represents the claims of a JWT.
type JWTClaims struct {
	jwt.Claims
	UserID string `json:"uid,omitempty"`
}

const (
	issuer = "https://github.com/cybersamx/authx"
)

var (
	ErrInvalidJWT = errors.New("invalid jwt")
	ErrExpiredJWT = errors.New("expired jwt")
)

func NewJWT(id, uid, secrets string, issueAt, expireAt time.Time) (string, error) {
	claims := JWTClaims{
		UserID: uid,
		Claims: jwt.Claims{
			Issuer:   issuer,
			Expiry:   jwt.NewNumericDate(expireAt),
			IssuedAt: jwt.NewNumericDate(issueAt),
			ID:       id,
		},
	}

	// For now, use symmetrical signing.
	signKey := jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       []byte(secrets),
	}
	opts := jose.SignerOptions{}
	opts.WithType("JWT")
	signer, err := jose.NewSigner(signKey, &opts)
	if err != nil {
		return "", err
	}

	signedJWT, err := jwt.
		Signed(signer).
		Claims(claims).
		CompactSerialize()

	return signedJWT, err
}

// ParseJWT parses the token string and validates the signature.
func ParseJWT(jwtStr, accessSecret string) (*JWTClaims, error) {
	signedJWT, err := jwt.ParseSigned(jwtStr)
	if err != nil {
		return nil, ErrInvalidJWT
	}

	claims := new(JWTClaims)
	if err = signedJWT.Claims([]byte(accessSecret), claims); err != nil {
		return nil, ErrInvalidJWT
	}

	err = claims.Validate(jwt.Expected{
		Issuer: issuer,
		Time:   time.Now(),
	})
	if err == jwt.ErrExpired {
		return nil, ErrExpiredJWT
	} else if err != nil {
		return nil, ErrInvalidJWT
	}

	return claims, nil
}

// UnsafeParseJWT parses the token string but no signature validation.
func UnsafeParseJWT(jwtStr string) (*JWTClaims, error) {
	signedJWT, err := jwt.ParseSigned(jwtStr)
	if err != nil {
		return nil, ErrInvalidJWT
	}

	claims := new(JWTClaims)
	if err := signedJWT.UnsafeClaimsWithoutVerification(claims); err != nil {
		return nil, ErrInvalidJWT
	}

	return claims, nil
}
