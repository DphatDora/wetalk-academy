package repository

import (
	"context"
	"wetalk-academy/internal/domain/model"
)

type QuizSubmissionRepository interface {
	CreateSubmission(ctx context.Context, submission *model.QuizSubmission) error
}
