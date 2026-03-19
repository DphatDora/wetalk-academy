package repository

import (
	"context"
	"wetalk-academy/internal/domain/model"
	"wetalk-academy/internal/domain/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TopicRepositoryImpl struct {
	collection *mongo.Collection
}

func NewTopicRepository(db *mongo.Database) repository.TopicRepository {
	return &TopicRepositoryImpl{
		collection: db.Collection(model.Topic{}.CollectionName()),
	}
}

func (r *TopicRepositoryImpl) CreateTopic(ctx context.Context, topic *model.Topic) error {
	_, err := r.collection.InsertOne(ctx, topic)
	return err
}

func (r *TopicRepositoryImpl) GetTopics(ctx context.Context, page, limit int) ([]*model.Topic, int64, error) {
	total, err := r.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var topics []*model.Topic
	if err := cursor.All(ctx, &topics); err != nil {
		return nil, 0, err
	}

	return topics, total, nil
}

func (r *TopicRepositoryImpl) GetTopicBySlug(ctx context.Context, slug string) (*model.Topic, error) {
	var topic *model.Topic
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&topic)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (r *TopicRepositoryImpl) GetTopicByID(ctx context.Context, id string) (*model.Topic, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var topic *model.Topic
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&topic)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (r *TopicRepositoryImpl) UpdateTopic(ctx context.Context, topic *model.Topic) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"slug": topic.Slug}, bson.M{"$set": topic})
	return err
}

func (r *TopicRepositoryImpl) DeleteTopic(ctx context.Context, slug string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"slug": slug})
	return err
}
