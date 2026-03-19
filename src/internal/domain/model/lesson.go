package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Lesson struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	TopicID    bson.ObjectID `bson:"topic_id"`
	Slug       string        `bson:"slug"`
	Title      string        `bson:"title"`
	OrderIndex int           `bson:"order_index"`
	CreatedAt  time.Time     `bson:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at"`
}

func (Lesson) CollectionName() string {
	return "lessons"
}
