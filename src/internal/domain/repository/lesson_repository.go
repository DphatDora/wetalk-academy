package repository

import (
	"context"
	"wetalk-academy/internal/domain/model"
)

type LessonRepository interface {
	CreateLesson(ctx context.Context, lesson *model.Lesson) error
	GetLessonBySlug(ctx context.Context, slug string) (*model.Lesson, error)
	GetLessonByID(ctx context.Context, id string) (*model.Lesson, error)
	GetLessonsByTopicID(ctx context.Context, topicID string, page, limit int) ([]*model.Lesson, int64, error)
	UpdateLesson(ctx context.Context, lesson *model.Lesson) error
	DeleteLesson(ctx context.Context, slug string) error
	CountLessonsByTopicID(ctx context.Context, topicID string) (int64, error)
}
