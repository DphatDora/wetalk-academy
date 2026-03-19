package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TopicAuthor struct {
	UserID uint64 `bson:"user_id"`
	Avatar string `bson:"avatar"`
	Name   string `bson:"name"`
}

type Topic struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Slug        string        `bson:"slug"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	Author      TopicAuthor   `bson:"author"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
}

func (Topic) CollectionName() string {
	return "topics"
}
