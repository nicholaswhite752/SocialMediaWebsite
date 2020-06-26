package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID       primitive.ObjectID `bson:"_id"`
	Caption string `bson:"caption"`
	Username string `bson:"username"`
	Filename string `bson:"filename"`
	Likes int `bson:"likes"`
	DatePosted time.Time `bson:"dateposted"`
	UserLikesIDS []primitive.ObjectID `bson:"userlikesids"`
}

type PostWithEncode struct {
	ID       primitive.ObjectID `bson:"_id"`
	Caption string `bson:"caption"`
	Username string `bson:"username"`
	Filename string `bson:"filename"`
	Likes int `bson:"likes"`
	DatePosted string `bson:"dateposted"`
	UserLikesIDS []primitive.ObjectID `bson:"userlikesids"`
	EncodedFile	string
	IsLiked string
	DeleteButton bool
}


type PostToDB struct {
	Username string `bson:"username"`
	Caption string `bson:"caption"`
	Filename string `bson:"filename"`
	Likes int `bson:"likes"`
	DatePosted time.Time `bson:"dateposted"`
	UserLikesIDS []primitive.ObjectID `bson:"userlikesids"`
}

type FSFile struct{
	ID       primitive.ObjectID `bson:"_id"`
	Length int `bson:"length"`
	ChunkSize int `bson:"chunksize"`
	UploadTime time.Time `bson:"uploadDate"`
	Filename string `bson:"filename"`
}