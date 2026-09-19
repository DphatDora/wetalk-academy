package data

import (
	"fmt"
	"time"
	"wetalk-academy/internal/domain/model"
)

var seedTime = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

var topicMediaURLs = map[string]string{
	"go-programming":             "https://go.dev/images/gophers/ladder.svg",
	"backend-go":                 "https://go.dev/images/gophers/pilot-bust.svg",
	"modern-js-ts":               "https://www.typescriptlang.org/images/branding/two-colors/ts-lettermark-white.svg",
	"react-engineering":          "https://react.dev/images/uwu.png",
	"python-automation-data":     "https://www.python.org/static/community_logos/python-logo-master-v3-TM.png",
	"devops-cicd":                "https://www.docker.com/wp-content/uploads/2022/03/Moby-logo.png",
	"database-engineering":       "https://webimages.mongodb.com/_com_assets/cms/kuy2cq2e4m3hy7u9d-DB_White.png",
	"system-design":              "https://images.unsplash.com/photo-1451187580459-43490279c0fa",
	"cloud-native-microservices": "https://kubernetes.io/images/kubernetes-horizontal-color.png",
	"web-security":               "https://owasp.org/assets/images/logo.png",
}

func All() []TopicDataset {
	specs := topicSpecs()
	result := make([]TopicDataset, 0, len(specs))
	for _, spec := range specs {
		result = append(result, buildTopic(spec))
	}
	return result
}

func Select(key string) ([]TopicDataset, error) {
	if key == "" {
		return All(), nil
	}
	for _, dataset := range All() {
		if dataset.Key == key {
			return []TopicDataset{dataset}, nil
		}
	}
	return nil, fmt.Errorf("unknown topic dataset %q", key)
}

func Keys() []string {
	all := All()
	keys := make([]string, len(all))
	for i := range all {
		keys[i] = all[i].Key
	}
	return keys
}

func buildTopic(spec TopicSpec) TopicDataset {
	topicID := stableID(DatasetVersion, spec.Key, "topic")
	if spec.Author.UserID == 0 {
		spec.Author = authorFor(spec.Key)
	}
	topic := model.Topic{
		ID:          topicID,
		Slug:        spec.Key,
		Title:       spec.Title,
		Description: spec.Description,
		Author:      spec.Author,
		CreatedAt:   seedTime,
		UpdatedAt:   seedTime,
	}

	dataset := TopicDataset{Key: spec.Key, Topic: topic}
	for index, lessonSpec := range spec.Lessons {
		lessonKey := fmt.Sprintf("%s-%02d", spec.Key, index+1)
		lessonID := stableID(DatasetVersion, lessonKey, "lesson")
		lesson := model.Lesson{
			ID:         lessonID,
			TopicID:    topicID,
			Slug:       lessonKey,
			Title:      lessonSpec.Title,
			OrderIndex: index + 1,
			CreatedAt:  seedTime,
			UpdatedAt:  seedTime,
		}
		dataset.Lessons = append(dataset.Lessons, lesson)
		dataset.Contents = append(dataset.Contents, buildContent(spec, lessonSpec, lesson))
		for quizIndex := 0; quizIndex < lessonSpec.QuizCount; quizIndex++ {
			dataset.Quizzes = append(dataset.Quizzes, buildQuiz(spec, lessonSpec, lesson, quizIndex))
		}
	}
	return dataset
}

func buildContent(topic TopicSpec, spec LessonSpec, lesson model.Lesson) model.Content {
	mediaURL := spec.MediaURL
	if mediaURL == "" {
		mediaURL = topicMediaURLs[topic.Key]
	}
	language := spec.Language
	if language == "" {
		language = topic.Language
	}

	sections := []model.ContentSection{
		{ID: 1, Type: "text", Content: fmt.Sprintf("# %s\n\n%s\n\n**Mục tiêu:** %s", lesson.Title, spec.Material.Overview, spec.Objective)},
		{ID: 2, Type: "text", Content: "## Kiến thức cốt lõi\n\n" + spec.Material.Concepts},
		{ID: 3, Type: "code", Language: language, Content: spec.Code},
		{ID: 4, Type: "text", Content: "## Tình huống thực tế và phân tích\n\n" + spec.Material.Scenario},
		{ID: 5, Type: "media", URL: mediaURL},
		{ID: 6, Type: "text", Content: "## Lỗi thường gặp và trade-off\n\n" + spec.Material.Pitfalls},
		{ID: 7, Type: "text", Content: "## Bài thực hành\n\n" + spec.Material.Exercise},
	}
	return model.Content{
		ID:        stableID(DatasetVersion, lesson.Slug, "content"),
		LessonID:  lesson.ID,
		Sections:  sections,
		CreatedAt: seedTime,
		UpdatedAt: seedTime,
	}
}

func buildQuiz(_ TopicSpec, spec LessonSpec, lesson model.Lesson, quizIndex int) model.Quiz {
	title := "Kiểm tra kiến thức"
	if quizIndex == 1 {
		title = "Bài tập tình huống"
	}
	questions := make([]model.QuizQuestion, 0, len(spec.Material.Quiz))
	for _, fact := range spec.Material.Quiz {
		question := fact.Question
		if quizIndex == 1 {
			question = "Trong một buổi review tình huống thực tế, " + lowerFirst(question)
		}
		questions = append(questions, model.QuizQuestion{
			Question: question, Point: 2,
			Options:       []string{fact.Correct, fact.Wrong[0], fact.Wrong[1], fact.Wrong[2]},
			CorrectAnswer: fact.Correct,
		})
	}
	return model.Quiz{
		ID:        stableID(DatasetVersion, lesson.Slug, "quiz", fmt.Sprint(quizIndex+1)),
		LessonID:  lesson.ID,
		Title:     fmt.Sprintf("%s: %s", lesson.Title, title),
		Questions: questions,
		TimeLimit: 480,
		CreatedAt: seedTime,
	}
}

func lowerFirst(value string) string {
	if value == "" {
		return value
	}
	first := value[0]
	if first >= 'A' && first <= 'Z' {
		first += 'a' - 'A'
	}
	return string(first) + value[1:]
}
