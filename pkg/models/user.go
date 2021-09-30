package models

type User struct {
	ID       string `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Salt     string `bson:"salt"`
}

func (u *User) RemoveSensitiveData() {
	u.Password = ""
	u.Salt = ""
}
