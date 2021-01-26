package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

// UserClaims represents the claims of a JWT.
type UserClaims struct {
	UserID string `json:"id"`
	jwt.Claims
}

const (
	issuer  = "Authx"
	subject = "Access token"
)

var (
	ErrInvalidJWT = errors.New("invalid jwt")
	ErrExpiredJWT = errors.New("expired jwt")
)

// GetBearerFromHeader gets the token token from the header string.
func GetBearerFromHeader(header string) (string, error) {
	if header == "" {
		return "", ErrMissingBearer
	}

	chunks := strings.Split(header, " ")
	if len(chunks) == 2 && chunks[0] == "Bearer" {
		return chunks[1], nil
	}

	return "", ErrInvalidBearer
}

func NewJWT(id, uid, secrets string, issueAt, expireAt time.Time) (string, error) {
	claims := UserClaims{
		UserID: uid,
		Claims: jwt.Claims{
			Issuer:    issuer,
			Subject:   subject,
			Expiry:    jwt.NewNumericDate(expireAt),
			NotBefore: jwt.NewNumericDate(issueAt),
			IssuedAt:  jwt.NewNumericDate(issueAt),
			ID:        id,
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

func ParseJWTFromRequest(req *http.Request, accessSecret string) (*UserClaims, error) {
	bearer, err := GetBearerFromHeader(req.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	return parseJWT(bearer, accessSecret)
}

func parseJWT(signedJWT, accessSecret string) (*UserClaims, error) {
	jwtToken, err := jwt.ParseSigned(signedJWT)
	if err != nil {
		return nil, ErrInvalidJWT
	}

	claims := new(UserClaims)
	if cerr := jwtToken.Claims([]byte(accessSecret), claims); cerr != nil {
		return nil, ErrInvalidJWT
	}

	err = claims.Validate(jwt.Expected{
		Issuer:  issuer,
		Subject: subject,
		Time:    time.Now(),
	})
	if err == jwt.ErrExpired {
		return nil, ErrExpiredJWT
	} else if err != nil {
		return nil, ErrInvalidJWT
	}

	return claims, nil
}
