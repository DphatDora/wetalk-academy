package runner

import (
	"context"
	"fmt"
	"time"
	"wetalk-academy/internal/domain/model"
	"wetalk-academy/seeders/data"
	"wetalk-academy/seeders/validation"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const manifestCollection = "seeder_runs"

type Manifest struct {
	ID         string          `bson:"_id"`
	Version    string          `bson:"version"`
	TopicIDs   []bson.ObjectID `bson:"topic_ids"`
	LessonIDs  []bson.ObjectID `bson:"lesson_ids"`
	ContentIDs []bson.ObjectID `bson:"content_ids"`
	QuizIDs    []bson.ObjectID `bson:"quiz_ids"`
	SeededAt   time.Time       `bson:"seeded_at"`
}

type Runner struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Runner {
	return &Runner{db: db}
}

func (r *Runner) Seed(ctx context.Context, datasets []data.TopicDataset) error {
	if err := validation.Validate(datasets); err != nil {
		return err
	}
	for _, dataset := range datasets {
		if err := r.seedTopic(ctx, dataset); err != nil {
			return fmt.Errorf("seed %s: %w", dataset.Key, err)
		}
	}
	return nil
}

func (r *Runner) seedTopic(ctx context.Context, dataset data.TopicDataset) error {
	if err := replace(ctx, r.db.Collection(model.Topic{}.CollectionName()), dataset.Topic.ID, dataset.Topic); err != nil {
		return err
	}
	for _, lesson := range dataset.Lessons {
		if err := replace(ctx, r.db.Collection(model.Lesson{}.CollectionName()), lesson.ID, lesson); err != nil {
			return err
		}
	}
	for _, content := range dataset.Contents {
		if err := replace(ctx, r.db.Collection(model.Content{}.CollectionName()), content.ID, content); err != nil {
			return err
		}
	}
	for _, quiz := range dataset.Quizzes {
		if err := replace(ctx, r.db.Collection(model.Quiz{}.CollectionName()), quiz.ID, quiz); err != nil {
			return err
		}
	}

	manifest := manifestFor(dataset)
	_, err := r.db.Collection(manifestCollection).ReplaceOne(
		ctx,
		bson.M{"_id": manifest.ID},
		manifest,
		options.Replace().SetUpsert(true),
	)
	return err
}

func replace(ctx context.Context, collection *mongo.Collection, id bson.ObjectID, document any) error {
	_, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, document, options.Replace().SetUpsert(true))
	return err
}

func manifestFor(dataset data.TopicDataset) Manifest {
	manifest := Manifest{
		ID:       dataset.Key,
		Version:  data.DatasetVersion,
		TopicIDs: []bson.ObjectID{dataset.Topic.ID},
		SeededAt: time.Now().UTC(),
	}
	for _, lesson := range dataset.Lessons {
		manifest.LessonIDs = append(manifest.LessonIDs, lesson.ID)
	}
	for _, content := range dataset.Contents {
		manifest.ContentIDs = append(manifest.ContentIDs, content.ID)
	}
	for _, quiz := range dataset.Quizzes {
		manifest.QuizIDs = append(manifest.QuizIDs, quiz.ID)
	}
	return manifest
}

func (r *Runner) Clear(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if err := r.clearTopic(ctx, key); err != nil {
			return fmt.Errorf("clear %s: %w", key, err)
		}
	}
	return nil
}

func (r *Runner) clearTopic(ctx context.Context, key string) error {
	var manifest Manifest
	err := r.db.Collection(manifestCollection).FindOne(ctx, bson.M{"_id": key}).Decode(&manifest)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		return err
	}

	operations := []struct {
		collection string
		ids        []bson.ObjectID
	}{
		{model.Quiz{}.CollectionName(), manifest.QuizIDs},
		{model.Content{}.CollectionName(), manifest.ContentIDs},
		{model.Lesson{}.CollectionName(), manifest.LessonIDs},
		{model.Topic{}.CollectionName(), manifest.TopicIDs},
	}
	for _, operation := range operations {
		if len(operation.ids) == 0 {
			continue
		}
		if _, err := r.db.Collection(operation.collection).DeleteMany(ctx, bson.M{"_id": bson.M{"$in": operation.ids}}); err != nil {
			return err
		}
	}
	_, err = r.db.Collection(manifestCollection).DeleteOne(ctx, bson.M{"_id": key})
	return err
}
