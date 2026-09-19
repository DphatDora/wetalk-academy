package runner

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"wetalk-academy/seeders/data"
	"wetalk-academy/seeders/migration"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestIntegrationSeedIsIdempotentAndClearPreservesForeignData(t *testing.T) {
	uri := os.Getenv("MONGO_TEST_URI")
	if uri == "" {
		t.Skip("MONGO_TEST_URI is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(fmt.Sprintf("wetalk_academy_seed_test_%d", time.Now().UnixNano()))
	defer db.Drop(context.Background())
	if err := migration.Up(ctx, db); err != nil {
		t.Fatal(err)
	}

	selected, err := data.Select("go-programming")
	if err != nil {
		t.Fatal(err)
	}
	r := New(db)
	if err := r.Seed(ctx, selected); err != nil {
		t.Fatal(err)
	}
	if err := r.Seed(ctx, selected); err != nil {
		t.Fatal(err)
	}
	assertCount(t, ctx, db, "topics", 1)
	assertCount(t, ctx, db, "lessons", 6)
	assertCount(t, ctx, db, "contents", 6)
	assertCount(t, ctx, db, "quizzes", 7)

	foreignID := bson.NewObjectID()
	_, err = db.Collection("topics").InsertOne(ctx, bson.M{
		"_id": foreignID, "slug": "foreign-topic", "title": "Foreign", "description": "must remain",
		"author":     bson.M{"user_id": int64(1), "avatar": "", "name": "test"},
		"created_at": time.Now(), "updated_at": time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Clear(ctx, []string{"go-programming"}); err != nil {
		t.Fatal(err)
	}
	assertCount(t, ctx, db, "topics", 1)
}

func assertCount(t *testing.T, ctx context.Context, db *mongo.Database, collection string, expected int64) {
	t.Helper()
	count, err := db.Collection(collection).CountDocuments(ctx, bson.M{})
	if err != nil {
		t.Fatal(err)
	}
	if count != expected {
		t.Fatalf("%s: expected %d documents, got %d", collection, expected, count)
	}
}
