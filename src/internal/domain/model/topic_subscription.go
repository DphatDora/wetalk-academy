package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TopicSubscription struct {
	ID           bson.ObjectID   `bson:"_id,omitempty"`
	TopicID      bson.ObjectID   `bson:"topic_id"`
	UserID       string          `bson:"user_id"`
	SubscribedAt time.Time       `bson:"subscribed_at"`
	LessonsDone  []bson.ObjectID `bson:"lessons_done"`
}

func (TopicSubscription) CollectionName() string {
	return "topic_subscriptions"
}
