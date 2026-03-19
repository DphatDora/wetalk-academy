package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ContentSection struct {
	ID       uint64 `bson:"id"`
	Type     string `bson:"type"`
	Content  string `bson:"content,omitempty"`
	Language string `bson:"language,omitempty"`
	URL      string `bson:"url,omitempty"`
}

type Content struct {
	ID        bson.ObjectID    `bson:"_id,omitempty"`
	LessonID  bson.ObjectID    `bson:"lesson_id"`
	Sections  []ContentSection `bson:"sections"`
	CreatedAt time.Time        `bson:"created_at"`
	UpdatedAt time.Time        `bson:"updated_at"`
}

func (Content) CollectionName() string {
	return "contents"
}
