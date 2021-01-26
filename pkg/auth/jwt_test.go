package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/square/go-jose/v3/jwt"
	"github.com/stretchr/testify/assert"

	"github.com/cybersamx/authx/pkg/utils"
)

// Notes:
// Since we need to regenerate the literal JWT values whenever we make changes to the claims or JWT logic,
// here are some info that can help.
//
// 12/1/2019 = 1575158400 (Unix epoch time)
// 1/1/2020  = 1577836800
// 1/1/2030  = 1893456000
//
// JWT header  = {"alg":"HS256","typ":"JWT"}
// JWT payload = {"exp":1893456000,"iat":1575158400,"id":"MyUserID","iss":"Authx","jti":"MyUniqueID","nbf":1575158400,"sub":"Access token"}

func Test_NewJWT(t *testing.T) {
	var tcases = []struct {
		description string
		expireAt    time.Time
		issueAt     time.Time
		id          string
		userID      string
		secrets     string
		expected    string
		pass        bool
	}{
		{
			description: "valid jwt",
			expireAt:    time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
			issueAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			id:          "MyUniqueID",
			userID:      "MyUserID",
			secrets:     "MySecrets",
			expected:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsImlhdCI6MTU3NzgzNjgwMCwiaWQiOiJNeVVzZXJJRCIsImlzcyI6IkF1dGh4IiwianRpIjoiTXlVbmlxdWVJRCIsIm5iZiI6MTU3NzgzNjgwMCwic3ViIjoiQWNjZXNzIHRva2VuIn0.1sEKBqQEzPjfo3VoNAkJRfLYHWlhD4Cxx0g2UM5p-NA", //nolint:lll
			pass:        true,
		},
	}

	for _, tcase := range tcases {
		val := tcase
		t.Run(val.description, func(t *testing.T) {
			jwtToken, err := NewJWT(val.id, val.userID, val.secrets, val.issueAt, val.expireAt)
			if val.pass {
				assert.NoError(t, err)
				assert.Equal(t, val.expected, jwtToken)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func Test_ParseJWT(t *testing.T) {
	var tcases = []struct {
		description string
		jwt         string
		expExpireAt time.Time
		expIssueAt  time.Time
		expID       string
		expUserID   string
		expIssuer   string
		expSubject  string
		secrets     string
		expError    error
	}{
		{
			description: "valid jwt",
			jwt:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsImlhdCI6MTU3NTE1ODQwMCwiaWQiOiJNeVVzZXJJRCIsImlzcyI6IkF1dGh4IiwianRpIjoiTXlVbmlxdWVJRCIsIm5iZiI6MTU3NTE1ODQwMCwic3ViIjoiQWNjZXNzIHRva2VuIn0.pJdYBs6gOSY8lcBGXrvfa4BAU4DI26z4nEbUs7PQtm0", //nolint:lll
			expExpireAt: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
			expIssueAt:  time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC),
			expID:       "MyUniqueID",
			expUserID:   "MyUserID",
			expIssuer:   "Authx",
			expSubject:  "Access token",
			secrets:     "MySecrets",
			expError:    nil,
		},
		{
			description: "expired jwt",
			jwt:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzUxNTg0MDAsImlhdCI6MTU3NTE1ODQwMCwiaWQiOiJNeVVzZXJJRCIsImlzcyI6IkF1dGh4IiwianRpIjoiTXlVbmlxdWVJRCIsIm5iZiI6MTU3NTE1ODQwMCwic3ViIjoiQWNjZXNzIHRva2VuIn0.3EmqBGvgnvDBCCRQ9oydsQJjgGSCcBvwCPJXwtltCHY", //nolint:lll
			expExpireAt: time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC),
			expIssueAt:  time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC),
			expID:       "MyUniqueID",
			expUserID:   "MyUserID",
			expIssuer:   "Authx",
			expSubject:  "Access token",
			secrets:     "MySecrets",
			expError:    ErrExpiredJWT,
		},
	}

	for _, tcase := range tcases {
		val := tcase
		t.Run(val.description, func(t *testing.T) {
			claims, err := parseJWT(val.jwt, val.secrets)
			if val.expError != nil {
				assert.Equal(t, val.expError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, val.expUserID, claims.UserID)
				assert.Equal(t, val.expID, claims.ID)
				assert.Equal(t, val.expSubject, claims.Subject)
				assert.Equal(t, val.expIssuer, claims.Issuer)
				assert.Equal(t, jwt.NewNumericDate(val.expIssueAt), claims.IssuedAt)
				assert.Equal(t, jwt.NewNumericDate(val.expExpireAt), claims.Expiry)
			}
		})
	}
}

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

			claims, err := parseJWT(jwtToken, val.secrets)
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
