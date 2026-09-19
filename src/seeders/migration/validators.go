package migration

import "go.mongodb.org/mongo-driver/v2/bson"

func validators() map[string]bson.M {
	objectID := bson.M{"bsonType": "objectId"}
	date := bson.M{"bsonType": "date"}
	stringValue := bson.M{"bsonType": "string"}
	positiveInt := bson.M{"bsonType": []string{"int", "long"}, "minimum": 1}

	return map[string]bson.M{
		"topics": schema(
			[]string{"slug", "title", "description", "author", "created_at", "updated_at"},
			bson.M{
				"_id": objectID, "slug": stringValue, "title": stringValue, "description": stringValue,
				"author":     bson.M{"bsonType": "object", "required": []string{"user_id", "avatar", "name"}},
				"created_at": date, "updated_at": date,
			},
		),
		"lessons": schema(
			[]string{"topic_id", "slug", "title", "order_index", "created_at", "updated_at"},
			bson.M{"_id": objectID, "topic_id": objectID, "slug": stringValue, "title": stringValue, "order_index": positiveInt, "created_at": date, "updated_at": date},
		),
		"contents": schema(
			[]string{"lesson_id", "sections", "created_at", "updated_at"},
			bson.M{
				"_id": objectID, "lesson_id": objectID, "created_at": date, "updated_at": date,
				"sections": bson.M{"bsonType": "array", "minItems": 1, "items": bson.M{
					"bsonType": "object", "required": []string{"id", "type"},
					"properties": bson.M{
						"id":      positiveInt,
						"type":    bson.M{"enum": []string{"text", "media", "code"}},
						"content": stringValue, "language": stringValue, "url": stringValue,
					},
					"oneOf": []bson.M{
						{"required": []string{"content"}, "properties": bson.M{"type": bson.M{"enum": []string{"text"}}}},
						{"required": []string{"url"}, "properties": bson.M{"type": bson.M{"enum": []string{"media"}}}},
						{"required": []string{"content", "language"}, "properties": bson.M{"type": bson.M{"enum": []string{"code"}}}},
					},
				}},
			},
		),
		"quizzes": schema(
			[]string{"lesson_id", "title", "questions", "time_limit", "created_at"},
			bson.M{
				"_id": objectID, "lesson_id": objectID, "title": stringValue, "time_limit": positiveInt, "created_at": date,
				"questions": bson.M{"bsonType": "array", "minItems": 1, "items": bson.M{
					"bsonType": "object", "required": []string{"question", "point", "options", "correct_answer"},
					"properties": bson.M{"question": stringValue, "point": positiveInt, "options": bson.M{"bsonType": "array", "minItems": 2, "items": stringValue}, "correct_answer": stringValue},
				}},
			},
		),
		"quiz_submissions": schema(
			[]string{"quiz_id", "user_id", "answers", "total_time", "total_score", "submitted_at"},
			bson.M{
				"_id": objectID, "quiz_id": objectID, "user_id": bson.M{"bsonType": []string{"int", "long"}},
				"answers": bson.M{"bsonType": "array", "items": stringValue}, "total_time": bson.M{"bsonType": []string{"int", "long"}, "minimum": 0},
				"total_score": bson.M{"bsonType": []string{"double", "int", "long"}, "minimum": 0, "maximum": 100}, "submitted_at": date,
			},
		),
		"topic_subscriptions": schema(
			[]string{"topic_id", "user_id", "subscribed_at", "lessons_done"},
			bson.M{
				"_id": objectID, "topic_id": objectID, "user_id": stringValue, "subscribed_at": date,
				"lessons_done": bson.M{"bsonType": "array", "items": objectID},
			},
		),
	}
}

func schema(required []string, properties bson.M) bson.M {
	return bson.M{"$jsonSchema": bson.M{
		"bsonType":             "object",
		"required":             required,
		"properties":           properties,
		"additionalProperties": true,
	}}
}
