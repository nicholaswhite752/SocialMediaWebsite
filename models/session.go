package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Session struct {
	UserID       primitive.ObjectID
	LastActivity time.Time
}
