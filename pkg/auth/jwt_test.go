package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/square/go-jose/v3/jwt"
	"github.com/stretchr/testify/assert"

	"github.com/cybersamx/authx/pkg/utils"
)

func Test_NewAndParseJWT(t *testing.T) {
	var tcases = []struct {
		description string
		expireAt    time.Time
		issueAt     time.Time
		id          string
		userID      string
		secrets     string
		expError    error
		toTamper    bool
	}{
		{
			description: "valid jwt",
			expireAt:    time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
			issueAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			id:          "MyUniqueID",
			userID:      "MyUserID",
			secrets:     "MySecrets",
		},
		{
			description: "valid jwt using other values",
			expireAt:    time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
			issueAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			id:          utils.GetRandString(8),
			userID:      utils.GetRandString(8),
			secrets:     utils.GetRandString(8),
		},
		{
			description: "invalid jwt",
			expireAt:    time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
			issueAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			id:          "MyUniqueID",
			userID:      "MyUserID",
			secrets:     "MySecrets",
			toTamper:    true,
			expError:    ErrInvalidJWT,
		},
		{
			description: "invalid jwt",
			expireAt:    time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC),
			issueAt:     time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC),
			id:          "MyUniqueID",
			userID:      "MyUserID",
			secrets:     "MySecrets",
			expError:    ErrExpiredJWT,
		},
	}

	for _, tcase := range tcases {
		val := tcase
		t.Run(val.description, func(t *testing.T) {
			jwtToken, err := NewJWT(val.id, val.userID, val.secrets, val.issueAt, val.expireAt)
			assert.NoError(t, err)
			assert.NotEmpty(t, jwtToken)

			if val.toTamper {
				tamper := utils.GetRandString(4)
				jwtToken = fmt.Sprintf("%s%s", jwtToken, tamper)
			}

			claims, err := ParseJWT(jwtToken, val.secrets)
			if val.expError != nil {
				assert.Equal(t, val.expError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, val.userID, claims.UserID)
				assert.Equal(t, val.id, claims.ID)
				assert.Equal(t, jwt.NewNumericDate(val.issueAt), claims.IssuedAt)
				assert.Equal(t, jwt.NewNumericDate(val.expireAt), claims.Expiry)
			}
		})
	}
}
