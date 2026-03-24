package response

import (
	"time"
	"wetalk-academy/internal/domain/model"
)

type ContentSectionResponse struct {
	ID       uint64 `json:"id"`
	Type     string `json:"type"`
	Content  string `json:"content,omitempty"`
	Language string `json:"language,omitempty"`
	URL      string `json:"url,omitempty"`
}

type ContentResponse struct {
	ID        string                   `json:"id"`
	LessonID  string                   `json:"lessonId"`
	Sections  []ContentSectionResponse `json:"sections"`
	CreatedAt time.Time                `json:"createdAt"`
	UpdatedAt time.Time                `json:"updatedAt"`
}

func NewContentResponse(content *model.Content) *ContentResponse {
	sections := make([]ContentSectionResponse, len(content.Sections))
	for i, section := range content.Sections {
		sections[i] = ContentSectionResponse{
			ID:       section.ID,
			Type:     section.Type,
			Content:  section.Content,
			Language: section.Language,
			URL:      section.URL,
		}
	}

	return &ContentResponse{
		ID:        content.ID.Hex(),
		LessonID:  content.LessonID.Hex(),
		Sections:  sections,
		CreatedAt: content.CreatedAt,
		UpdatedAt: content.UpdatedAt,
	}
}
