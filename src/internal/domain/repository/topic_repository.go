package repository

import (
	"context"
	"wetalk-academy/internal/domain/model"
)

type TopicRepository interface {
	CreateTopic(ctx context.Context, topic *model.Topic) error
	GetTopics(ctx context.Context, page, limit int) ([]*model.Topic, int64, error)
	GetTopicBySlug(ctx context.Context, slug string) (*model.Topic, error)
	GetTopicByID(ctx context.Context, id string) (*model.Topic, error)
	UpdateTopic(ctx context.Context, topic *model.Topic) error
	DeleteTopic(ctx context.Context, slug string) error
}
