package migration

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type indexDefinition struct {
	collection string
	name       string
	keys       bson.D
	unique     bool
}

var indexes = []indexDefinition{
	{"topics", "seed_topics_slug_unique", bson.D{{Key: "slug", Value: 1}}, true},
	{"lessons", "seed_lessons_slug_unique", bson.D{{Key: "slug", Value: 1}}, true},
	{"lessons", "seed_lessons_topic_order_unique", bson.D{{Key: "topic_id", Value: 1}, {Key: "order_index", Value: 1}}, true},
	{"lessons", "seed_lessons_topic_lookup", bson.D{{Key: "topic_id", Value: 1}}, false},
	{"contents", "seed_contents_lesson_unique", bson.D{{Key: "lesson_id", Value: 1}}, true},
	{"quizzes", "seed_quizzes_lesson_lookup", bson.D{{Key: "lesson_id", Value: 1}}, false},
	{"quiz_submissions", "seed_quiz_submissions_quiz_lookup", bson.D{{Key: "quiz_id", Value: 1}}, false},
	{"topic_subscriptions", "seed_topic_subscriptions_topic_user_unique", bson.D{{Key: "topic_id", Value: 1}, {Key: "user_id", Value: 1}}, true},
	{"topic_subscriptions", "seed_topic_subscriptions_topic_lookup", bson.D{{Key: "topic_id", Value: 1}}, false},
}

func Up(ctx context.Context, db *mongo.Database) error {
	for collection, validator := range validators() {
		if err := applyValidator(ctx, db, collection, validator); err != nil {
			return fmt.Errorf("validator %s: %w", collection, err)
		}
	}
	for _, definition := range indexes {
		if definition.unique {
			if err := ensureNoDuplicates(ctx, db.Collection(definition.collection), definition.keys); err != nil {
				return fmt.Errorf("index %s: %w", definition.name, err)
			}
		}
		model := mongo.IndexModel{
			Keys: definition.keys,
			Options: options.Index().
				SetName(definition.name).
				SetUnique(definition.unique),
		}
		if _, err := db.Collection(definition.collection).Indexes().CreateOne(ctx, model); err != nil {
			return fmt.Errorf("create index %s: %w", definition.name, err)
		}
	}
	return nil
}

func Down(ctx context.Context, db *mongo.Database) error {
	for i := len(indexes) - 1; i >= 0; i-- {
		definition := indexes[i]
		err := db.Collection(definition.collection).Indexes().DropOne(ctx, definition.name)
		if err != nil && !isMissingNamespaceOrIndex(err) {
			return fmt.Errorf("drop index %s: %w", definition.name, err)
		}
	}
	for collection := range validators() {
		command := bson.D{
			{Key: "collMod", Value: collection},
			{Key: "validator", Value: bson.M{}},
			{Key: "validationLevel", Value: "off"},
		}
		if err := db.RunCommand(ctx, command).Err(); err != nil && !isMissingNamespaceOrIndex(err) {
			return fmt.Errorf("remove validator %s: %w", collection, err)
		}
	}
	return nil
}

func isMissingNamespaceOrIndex(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "index not found") ||
		strings.Contains(message, "namespace not found") ||
		strings.Contains(message, "namespacenotfound")
}

func ensureNoDuplicates(ctx context.Context, collection *mongo.Collection, keys bson.D) error {
	groupID := bson.D{}
	for _, key := range keys {
		groupID = append(groupID, bson.E{Key: key.Key, Value: "$" + key.Key})
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: groupID},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "count", Value: bson.D{{Key: "$gt", Value: 1}}}}}},
		bson.D{{Key: "$limit", Value: 1}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		var duplicate bson.M
		if err := cursor.Decode(&duplicate); err != nil {
			return err
		}
		return fmt.Errorf("duplicate data exists: %v", duplicate["_id"])
	}
	return cursor.Err()
}

func applyValidator(ctx context.Context, db *mongo.Database, collection string, validator bson.M) error {
	names, err := db.ListCollectionNames(ctx, bson.M{"name": collection})
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return db.CreateCollection(ctx, collection, options.CreateCollection().
			SetValidator(validator).
			SetValidationLevel("strict").
			SetValidationAction("error"))
	}
	command := bson.D{
		{Key: "collMod", Value: collection},
		{Key: "validator", Value: validator},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}
	return db.RunCommand(ctx, command).Err()
}
