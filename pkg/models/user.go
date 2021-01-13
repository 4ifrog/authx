package models

type User struct {
	ID       string `json:"id"       bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Salt     string `json:"salt"     bson:"salt"`
}
