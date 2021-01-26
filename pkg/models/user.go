package models

type User struct {
	ID       string `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Salt     string `bson:"salt"`
}
