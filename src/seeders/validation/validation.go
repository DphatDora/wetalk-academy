package validation

import (
	"fmt"
	"strings"
	"wetalk-academy/internal/domain/model"
	"wetalk-academy/seeders/data"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func Validate(datasets []data.TopicDataset) error {
	var problems []string
	topicIDs := map[bson.ObjectID]bool{}
	lessonIDs := map[bson.ObjectID]bool{}
	slugs := map[string]string{}
	authors := map[uint64]bool{}

	for _, dataset := range datasets {
		if dataset.Key == "" || dataset.Topic.Slug != dataset.Key {
			problems = append(problems, fmt.Sprintf("dataset %q: key must match topic slug", dataset.Key))
		}
		checkSlug(&problems, slugs, dataset.Topic.Slug, "topic")
		topicIDs[dataset.Topic.ID] = true
		if dataset.Topic.Author.UserID == 0 || strings.TrimSpace(dataset.Topic.Author.Name) == "" || !strings.HasPrefix(dataset.Topic.Author.Avatar, "https://") {
			problems = append(problems, fmt.Sprintf("%s: topic author must have ID, name, and HTTPS avatar", dataset.Key))
		}
		authors[dataset.Topic.Author.UserID] = true
		orderIndexes := map[int]bool{}
		for _, lesson := range dataset.Lessons {
			checkSlug(&problems, slugs, lesson.Slug, "lesson")
			if lesson.TopicID != dataset.Topic.ID {
				problems = append(problems, fmt.Sprintf("%s: invalid topic reference", lesson.Slug))
			}
			if lesson.OrderIndex < 1 || orderIndexes[lesson.OrderIndex] {
				problems = append(problems, fmt.Sprintf("%s: duplicate or invalid order_index %d", lesson.Slug, lesson.OrderIndex))
			}
			orderIndexes[lesson.OrderIndex] = true
			lessonIDs[lesson.ID] = true
		}
		for index := 1; index <= len(dataset.Lessons); index++ {
			if !orderIndexes[index] {
				problems = append(problems, fmt.Sprintf("%s: missing order_index %d", dataset.Key, index))
			}
		}
	}

	contentByLesson := map[bson.ObjectID]bool{}
	contentFingerprints := map[string]string{}
	quizQuestions := map[string]string{}
	for _, dataset := range datasets {
		for _, content := range dataset.Contents {
			if !lessonIDs[content.LessonID] {
				problems = append(problems, fmt.Sprintf("content %s: unknown lesson", content.ID.Hex()))
			}
			if contentByLesson[content.LessonID] {
				problems = append(problems, fmt.Sprintf("lesson %s: multiple content documents", content.LessonID.Hex()))
			}
			contentByLesson[content.LessonID] = true
			validateSections(&problems, content.LessonID.Hex(), content.Sections)
			fingerprint := contentFingerprint(content.Sections)
			if previous, exists := contentFingerprints[fingerprint]; exists {
				problems = append(problems, fmt.Sprintf("lessons %s and %s have duplicated content", previous, content.LessonID.Hex()))
			}
			contentFingerprints[fingerprint] = content.LessonID.Hex()
		}
		for _, quiz := range dataset.Quizzes {
			if !lessonIDs[quiz.LessonID] {
				problems = append(problems, fmt.Sprintf("quiz %s: unknown lesson", quiz.ID.Hex()))
			}
			if quiz.TimeLimit < 1 || len(quiz.Questions) == 0 {
				problems = append(problems, fmt.Sprintf("quiz %s: questions and positive time_limit are required", quiz.ID.Hex()))
			}
			for i, question := range quiz.Questions {
				if strings.TrimSpace(question.Question) == "" || question.Point < 1 || len(question.Options) < 2 {
					problems = append(problems, fmt.Sprintf("quiz %s question %d: invalid question", quiz.ID.Hex(), i+1))
				}
				if !contains(question.Options, question.CorrectAnswer) {
					problems = append(problems, fmt.Sprintf("quiz %s question %d: correct answer is not an option", quiz.ID.Hex(), i+1))
				}
				normalized := strings.ToLower(strings.TrimSpace(question.Question))
				if previous, exists := quizQuestions[normalized]; exists {
					problems = append(problems, fmt.Sprintf("quiz %s question %d duplicates question from %s", quiz.ID.Hex(), i+1, previous))
				}
				quizQuestions[normalized] = quiz.ID.Hex()
			}
		}
	}
	for lessonID := range lessonIDs {
		if !contentByLesson[lessonID] {
			problems = append(problems, fmt.Sprintf("lesson %s: content is required", lessonID.Hex()))
		}
	}
	if len(topicIDs) != len(datasets) {
		problems = append(problems, "duplicate topic IDs detected")
	}
	if len(datasets) > 1 && len(authors) < 3 {
		problems = append(problems, "multi-topic dataset must use at least 3 different authors")
	}
	if len(problems) > 0 {
		return fmt.Errorf("dataset validation failed:\n- %s", strings.Join(problems, "\n- "))
	}
	return nil
}

func checkSlug(problems *[]string, seen map[string]string, slug, kind string) {
	if strings.TrimSpace(slug) == "" {
		*problems = append(*problems, kind+": empty slug")
		return
	}
	if previous, exists := seen[slug]; exists {
		*problems = append(*problems, fmt.Sprintf("%s slug %q duplicates %s", kind, slug, previous))
	}
	seen[slug] = kind
}

func validateSections(problems *[]string, lesson string, sections []model.ContentSection) {
	if len(sections) < 6 {
		*problems = append(*problems, fmt.Sprintf("lesson %s: at least 6 content sections are required", lesson))
		return
	}
	totalTextLength := 0
	hasCode := false
	hasMedia := false
	for i, section := range sections {
		if section.ID != uint64(i+1) {
			*problems = append(*problems, fmt.Sprintf("lesson %s: section IDs must be sequential", lesson))
		}
		switch section.Type {
		case "text":
			if strings.TrimSpace(section.Content) == "" {
				*problems = append(*problems, fmt.Sprintf("lesson %s section %d: text content is required", lesson, i+1))
			}
			totalTextLength += len([]rune(section.Content))
		case "code":
			hasCode = true
			if strings.TrimSpace(section.Content) == "" || strings.TrimSpace(section.Language) == "" {
				*problems = append(*problems, fmt.Sprintf("lesson %s section %d: code and language are required", lesson, i+1))
			}
		case "media":
			hasMedia = true
			if !strings.HasPrefix(section.URL, "https://") {
				*problems = append(*problems, fmt.Sprintf("lesson %s section %d: HTTPS media URL is required", lesson, i+1))
			}
		default:
			*problems = append(*problems, fmt.Sprintf("lesson %s section %d: unsupported type %q", lesson, i+1, section.Type))
		}
	}
	if totalTextLength < 1200 {
		*problems = append(*problems, fmt.Sprintf("lesson %s: text content is too short (%d characters)", lesson, totalTextLength))
	}
	if !hasCode || !hasMedia {
		*problems = append(*problems, fmt.Sprintf("lesson %s: code and media sections are required", lesson))
	}
}

func contentFingerprint(sections []model.ContentSection) string {
	var builder strings.Builder
	for _, section := range sections {
		builder.WriteString(section.Type)
		builder.WriteString("|")
		builder.WriteString(section.Content)
		builder.WriteString("|")
		builder.WriteString(section.URL)
	}
	return builder.String()
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
