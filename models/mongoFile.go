package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MongoFile struct {
	ID       primitive.ObjectID `bson:"_id"`
	FileLength int `bson:"length"`
	FileChunkSize int `bson:"chunkSize"`
	FileUploadTime time.Time `bson:"uploadDate"`
	FileFilename string `bson:"filename"`
}
