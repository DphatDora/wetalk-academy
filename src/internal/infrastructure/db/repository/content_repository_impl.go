package repository

import (
	"context"
	"wetalk-academy/internal/domain/model"
	"wetalk-academy/internal/domain/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ContentRepositoryImpl struct {
	collection *mongo.Collection
}

func NewContentRepository(db *mongo.Database) repository.ContentRepository {
	return &ContentRepositoryImpl{
		collection: db.Collection(model.Content{}.CollectionName()),
	}
}

func (r *ContentRepositoryImpl) CreateContent(ctx context.Context, content *model.Content) error {
	_, err := r.collection.InsertOne(ctx, content)
	return err
}

func (r *ContentRepositoryImpl) GetContentByLessonID(ctx context.Context, lessonID string) (*model.Content, error) {
	objectID, err := bson.ObjectIDFromHex(lessonID)
	if err != nil {
		return nil, err
	}

	var content *model.Content
	err = r.collection.FindOne(ctx, bson.M{"lesson_id": objectID}).Decode(&content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (r *ContentRepositoryImpl) UpdateContent(ctx context.Context, content *model.Content) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": content.ID},
		bson.M{"$set": content},
	)
	return err
}

func (r *ContentRepositoryImpl) DeleteContentByLessonID(ctx context.Context, lessonID string) error {
	objectID, err := bson.ObjectIDFromHex(lessonID)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"lesson_id": objectID})
	return err
}

func (r *ContentRepositoryImpl) ContentExistsByLessonID(ctx context.Context, lessonID string) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(lessonID)
	if err != nil {
		return false, err
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"lesson_id": objectID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
