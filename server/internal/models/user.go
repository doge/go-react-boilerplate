package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID       bson.ObjectID `json:"_id" bson:"_id"`
	Username string        `json:"username" bson:"username"`
	Email    string        `json:"email" bson:"email"`
	Password string        `json:"-" bson:"password"`
}
