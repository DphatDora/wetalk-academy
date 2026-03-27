package repository

import (
	"context"
	"wetalk-academy/internal/domain/model"
	"wetalk-academy/internal/domain/repository"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type QuizSubmissionRepositoryImpl struct {
	collection *mongo.Collection
}

func NewQuizSubmissionRepository(db *mongo.Database) repository.QuizSubmissionRepository {
	return &QuizSubmissionRepositoryImpl{
		collection: db.Collection(model.QuizSubmission{}.CollectionName()),
	}
}

func (r *QuizSubmissionRepositoryImpl) CreateSubmission(ctx context.Context, submission *model.QuizSubmission) error {
	_, err := r.collection.InsertOne(ctx, submission)
	return err
}
