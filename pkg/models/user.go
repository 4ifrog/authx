package models

type User struct {
	ID       string `bson:"_id" form:"-"`
	Username string `bson:"username" form:"username" binding:"required"`
	Password string `bson:"password" form:"password" binding:"required"`
	Salt     string `bson:"salt" form:"-"`
}
