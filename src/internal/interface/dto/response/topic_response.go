package response

import (
	"time"
	"wetalk-academy/internal/domain/model"
)

type TopicAuthorResponse struct {
	UserID uint64 `json:"userId"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
}

type TopicResponse struct {
	ID          string              `json:"id"`
	Slug        string              `json:"slug"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Author      TopicAuthorResponse `json:"author"`
	CreatedAt   time.Time           `json:"createdAt"`
}

func NewTopicResponse(topic *model.Topic) *TopicResponse {
	return &TopicResponse{
		ID:          topic.ID.Hex(),
		Slug:        topic.Slug,
		Title:       topic.Title,
		Description: topic.Description,
		Author: TopicAuthorResponse{
			UserID: topic.Author.UserID,
			Avatar: topic.Author.Avatar,
			Name:   topic.Author.Name,
		},
		CreatedAt: topic.CreatedAt,
	}
}

func NewTopicListResponse(topics []*model.Topic) []*TopicResponse {
	result := make([]*TopicResponse, len(topics))
	for i, topic := range topics {
		result[i] = NewTopicResponse(topic)
	}
	return result
}
