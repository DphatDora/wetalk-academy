package response

import (
	"time"
	"wetalk-academy/internal/domain/model"
)

type LessonResponse struct {
	ID         string    `json:"id"`
	TopicID    string    `json:"topicId"`
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
	OrderIndex int       `json:"orderIndex"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func NewLessonResponse(lesson *model.Lesson) *LessonResponse {
	return &LessonResponse{
		ID:         lesson.ID.Hex(),
		TopicID:    lesson.TopicID.Hex(),
		Slug:       lesson.Slug,
		Title:      lesson.Title,
		OrderIndex: lesson.OrderIndex,
		CreatedAt:  lesson.CreatedAt,
		UpdatedAt:  lesson.UpdatedAt,
	}
}

func NewLessonListResponse(lessons []*model.Lesson) []*LessonResponse {
	result := make([]*LessonResponse, len(lessons))
	for i, lesson := range lessons {
		result[i] = NewLessonResponse(lesson)
	}
	return result
}
