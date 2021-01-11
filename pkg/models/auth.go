package models

import (
	"time"
)

type AccessToken struct {
	ID       string
	Value    string
	ExpireAt time.Time
}

type RefreshToken struct {
	ID       string
	Value    string
	ExpireAt time.Time
}
