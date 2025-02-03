package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Content   string             `bson:"content" json:"content"`
	Author    PostAuthor         `bson:"author" json:"author"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type PostAuthor struct {
	Name     string `bson:"name" json:"name"`
	Username string `bson:"username" json:"username"`
}

type CreatePostInput struct {
	Content string `json:"content" binding:"required"`
}
