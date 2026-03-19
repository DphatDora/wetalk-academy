package repository

import (
	"context"
	"wetalk-academy/internal/domain/model"
	"wetalk-academy/internal/domain/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LessonRepositoryImpl struct {
	collection *mongo.Collection
}

func NewLessonRepository(db *mongo.Database) repository.LessonRepository {
	return &LessonRepositoryImpl{
		collection: db.Collection(model.Lesson{}.CollectionName()),
	}
}

func (r *LessonRepositoryImpl) CreateLesson(ctx context.Context, lesson *model.Lesson) error {
	_, err := r.collection.InsertOne(ctx, lesson)
	return err
}

func (r *LessonRepositoryImpl) GetLessonBySlug(ctx context.Context, slug string) (*model.Lesson, error) {
	var lesson *model.Lesson
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&lesson)
	if err != nil {
		return nil, err
	}
	return lesson, nil
}

func (r *LessonRepositoryImpl) GetLessonsByTopicID(ctx context.Context, topicID string, page, limit int) ([]*model.Lesson, int64, error) {
	objectID, err := bson.ObjectIDFromHex(topicID)
	if err != nil {
		return nil, 0, err
	}

	filter := bson.M{"topic_id": objectID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	opts := options.Find().
		SetSort(bson.D{{Key: "order_index", Value: 1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var lessons []*model.Lesson
	if err := cursor.All(ctx, &lessons); err != nil {
		return nil, 0, err
	}

	return lessons, total, nil
}

func (r *LessonRepositoryImpl) UpdateLesson(ctx context.Context, lesson *model.Lesson) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"slug": lesson.Slug}, bson.M{"$set": lesson})
	return err
}

func (r *LessonRepositoryImpl) DeleteLesson(ctx context.Context, slug string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"slug": slug})
	return err
}

func (r *LessonRepositoryImpl) CountLessonsByTopicID(ctx context.Context, topicID string) (int64, error) {
	objectID, err := bson.ObjectIDFromHex(topicID)
	if err != nil {
		return 0, err
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"topic_id": objectID})
	return count, err
}
