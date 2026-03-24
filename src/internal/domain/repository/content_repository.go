package repository

import (
	"context"
	"wetalk-academy/internal/domain/model"
)

type ContentRepository interface {
	CreateContent(ctx context.Context, content *model.Content) error
	GetContentByLessonID(ctx context.Context, lessonID string) (*model.Content, error)
	UpdateContent(ctx context.Context, content *model.Content) error
	DeleteContentByLessonID(ctx context.Context, lessonID string) error
	ContentExistsByLessonID(ctx context.Context, lessonID string) (bool, error)
}
