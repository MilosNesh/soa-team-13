package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	ImageUrl  string             `bson:"image_url,omitempty" json:"image_url,omitempty"`
	Likes     []string           `bson:"likes,omitempty" json:"likes,omitempty"`
	UserId    string             `bson:"user_id" json:"userId"`
}
