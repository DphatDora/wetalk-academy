package repository

import (
	"context"
	"wetalk-academy/internal/domain/model"
	"wetalk-academy/internal/domain/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type QuizRepositoryImpl struct {
	collection *mongo.Collection
}

func NewQuizRepository(db *mongo.Database) repository.QuizRepository {
	return &QuizRepositoryImpl{
		collection: db.Collection(model.Quiz{}.CollectionName()),
	}
}

func (r *QuizRepositoryImpl) CreateQuiz(ctx context.Context, quiz *model.Quiz) error {
	_, err := r.collection.InsertOne(ctx, quiz)
	return err
}

func (r *QuizRepositoryImpl) GetQuizByID(ctx context.Context, id string) (*model.Quiz, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var quiz model.Quiz
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&quiz)
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *QuizRepositoryImpl) GetQuizzesByLessonID(ctx context.Context, lessonID string) ([]*model.Quiz, error) {
	objectID, err := bson.ObjectIDFromHex(lessonID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"lesson_id": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var quizzes []*model.Quiz
	if err := cursor.All(ctx, &quizzes); err != nil {
		return nil, err
	}
	return quizzes, nil
}


func (r *QuizRepositoryImpl) UpdateQuiz(ctx context.Context, quiz *model.Quiz) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": quiz.ID}, bson.M{"$set": quiz})
	return err
}

func (r *QuizRepositoryImpl) DeleteQuiz(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
