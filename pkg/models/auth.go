package models

import (
	"time"
)

type AccessToken struct {
	ID       string    `json:"id"        bson:"_id"`
	Value    string    `json:"value"     bson:"value"`
	UserID   string    `json:"user_id"   bson:"userID"`
	ExpireAt time.Time `json:"expire_at" bson:"expireAt"`
}

type RefreshToken struct {
	ID       string    `json:"id"        bson:"_id"`
	Value    string    `json:"value"     bson:"value"`
	UserID   string    `json:"user_id"   bson:"userID"`
	ExpireAt time.Time `json:"expire_at" bson:"expireAt"`
}
