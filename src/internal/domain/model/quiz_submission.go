package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type QuizSubmission struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	QuizID      bson.ObjectID `bson:"quiz_id"`
	UserID      string        `bson:"user_id"`
	Answers     []string      `bson:"answers"`
	TotalTime   int           `bson:"total_time"`
	TotalScore  float64       `bson:"total_score"`
	SubmittedAt time.Time     `bson:"submitted_at"`
}

func (QuizSubmission) CollectionName() string {
	return "quiz_submissions"
}
