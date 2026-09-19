package data

import "wetalk-academy/internal/domain/model"

const DatasetVersion = "academy-v1"

type TopicDataset struct {
	Key      string
	Topic    model.Topic
	Lessons  []model.Lesson
	Contents []model.Content
	Quizzes  []model.Quiz
}

type LessonSpec struct {
	Title     string
	Objective string
	Language  string
	Code      string
	MediaURL  string
	QuizCount int
	Material  LessonMaterial
}

type TopicSpec struct {
	Key         string
	Title       string
	Description string
	Language    string
	Author      model.TopicAuthor
	Lessons     []LessonSpec
}

type LessonMaterial struct {
	Overview string
	Concepts string
	Scenario string
	Pitfalls string
	Exercise string
	Quiz     []QuizFact
}

type QuizFact struct {
	Question string
	Correct  string
	Wrong    [3]string
}
