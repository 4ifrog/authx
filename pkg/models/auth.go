package models

import (
	"time"
)

type AccessToken struct {
	ID       string    `bson:"_id"`
	Value    string    `bson:"value"`
	ExpireAt time.Time `bson:"expireAt"`
}

type RefreshToken struct {
	ID       string    `bson:"_id"`
	Value    string    `bson:"value"`
	ExpireAt time.Time `bson:"expireAt"`
}
