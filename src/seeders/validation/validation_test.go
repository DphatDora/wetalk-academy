package validation

import (
	"strings"
	"testing"
	"wetalk-academy/seeders/data"
)

func TestValidateAll(t *testing.T) {
	if err := Validate(data.All()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateRejectsInvalidQuizAnswer(t *testing.T) {
	datasets := data.All()
	datasets[0].Quizzes[0].Questions[0].CorrectAnswer = "not-an-option"
	err := Validate(datasets)
	if err == nil || !strings.Contains(err.Error(), "correct answer is not an option") {
		t.Fatalf("expected invalid answer error, got %v", err)
	}
}

func TestValidateRejectsDuplicateOrder(t *testing.T) {
	datasets := data.All()
	datasets[0].Lessons[1].OrderIndex = datasets[0].Lessons[0].OrderIndex
	err := Validate(datasets)
	if err == nil || !strings.Contains(err.Error(), "duplicate or invalid order_index") {
		t.Fatalf("expected order error, got %v", err)
	}
}

func TestValidateRejectsInvalidSection(t *testing.T) {
	datasets := data.All()
	datasets[0].Contents[0].Sections[0].Content = ""
	err := Validate(datasets)
	if err == nil || !strings.Contains(err.Error(), "text content is required") {
		t.Fatalf("expected content error, got %v", err)
	}
}
