package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type QuizQuestion struct {
	Question      string   `bson:"question"`
	Point         int      `bson:"point"`
	Options       []string `bson:"options"`
	CorrectAnswer string   `bson:"correct_answer"`
}

type Quiz struct {
	ID        bson.ObjectID  `bson:"_id,omitempty"`
	LessonID  bson.ObjectID  `bson:"lesson_id"`
	Title     string         `bson:"title"`
	Questions []QuizQuestion `bson:"questions"`
	TimeLimit int            `bson:"time_limit"`
	CreatedAt time.Time      `bson:"created_at"`
}

func (Quiz) CollectionName() string {
	return "quizzes"
}
