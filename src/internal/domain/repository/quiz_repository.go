package repository

import (
	"context"
	"wetalk-academy/internal/domain/model"
)

type QuizRepository interface {
	CreateQuiz(ctx context.Context, quiz *model.Quiz) error
	GetQuizByID(ctx context.Context, id string) (*model.Quiz, error)
	GetQuizzesByLessonID(ctx context.Context, lessonID string) ([]*model.Quiz, error)
	UpdateQuiz(ctx context.Context, quiz *model.Quiz) error
	DeleteQuiz(ctx context.Context, id string) error
}
