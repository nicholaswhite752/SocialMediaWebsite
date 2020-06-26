package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	PostsIDS []primitive.ObjectID `bson:"postids"`
	LikesIDS []primitive.ObjectID `bson:"likesids"`
}

type UserToDB struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	PostsIDS []primitive.ObjectID `bson:"postids"`
	LikesIDS []primitive.ObjectID `bson:"likesids"`
}

