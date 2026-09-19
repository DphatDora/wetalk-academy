package data

import "testing"

func TestDatasetShapeAndStableIDs(t *testing.T) {
	first := All()
	second := All()
	if len(first) != 10 {
		t.Fatalf("expected 10 topics, got %d", len(first))
	}
	var lessons, quizzes int
	for i := range first {
		if len(first[i].Lessons) != 6 || len(first[i].Contents) != 6 {
			t.Fatalf("%s: expected 6 lessons and contents", first[i].Key)
		}
		if first[i].Topic.ID != second[i].Topic.ID || first[i].Lessons[0].ID != second[i].Lessons[0].ID {
			t.Fatalf("%s: IDs are not stable", first[i].Key)
		}
		lessons += len(first[i].Lessons)
		quizzes += len(first[i].Quizzes)
	}
	if lessons != 60 || quizzes < 50 {
		t.Fatalf("unexpected totals lessons=%d quizzes=%d", lessons, quizzes)
	}
}
