package service

import (
	"context"
	"fmt"
	"log"
	"time"
	"wetalk-academy/internal/domain/model"
	"wetalk-academy/internal/domain/repository"
	"wetalk-academy/internal/interface/dto/request"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/package/util"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TopicService struct {
	topicRepo  repository.TopicRepository
	lessonRepo repository.LessonRepository
}

func NewTopicService(topicRepo repository.TopicRepository, lessonRepo repository.LessonRepository) *TopicService {
	return &TopicService{
		topicRepo:  topicRepo,
		lessonRepo: lessonRepo,
	}
}

func (s *TopicService) CreateTopic(ctx context.Context, userId uint64, req *request.CreateTopicRequest) error {
	topic := &model.Topic{
		Slug:        util.GenerateSlug(req.Title),
		Title:       req.Title,
		Description: req.Description,
		Author: model.TopicAuthor{
			UserID: userId,
			Avatar: req.Author.Avatar,
			Name:   req.Author.Name,
		},
		CreatedAt: time.Now(),
	}

	if err := s.topicRepo.CreateTopic(ctx, topic); err != nil {
		log.Printf("[Err] Error creating topic in TopicService.CreateTopic: %v", err)
		return fmt.Errorf("failed to create topic: %w", err)
	}

	return nil
}

func (s *TopicService) GetTopics(ctx context.Context, page, limit int) ([]*response.TopicResponse, int64, error) {
	topics, total, err := s.topicRepo.GetTopics(ctx, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting topics in TopicService.GetTopics: %v", err)
		return nil, 0, fmt.Errorf("failed to get topics: %w", err)
	}

	topicResponses := make([]*response.TopicResponse, len(topics))
	for i, topic := range topics {
		topicResponses[i] = response.NewTopicResponse(topic)
	}

	return topicResponses, total, nil
}

func (s *TopicService) GetTopicBySlug(ctx context.Context, slug string) (*response.TopicResponse, error) {
	topic, err := s.topicRepo.GetTopicBySlug(ctx, slug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("topic not found")
		}
		log.Printf("[Err] Error getting topic in TopicService.GetTopicBySlug: %v", err)
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	return response.NewTopicResponse(topic), nil
}

func (s *TopicService) UpdateTopic(ctx context.Context, slug string, userId uint64, req *request.UpdateTopicRequest) error {
	existingTopic, err := s.topicRepo.GetTopicBySlug(ctx, slug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("topic not found")
		}
		log.Printf("[Err] Error getting topic in TopicService.UpdateTopic: %v", err)
		return fmt.Errorf("failed to get topic: %w", err)
	}

	if existingTopic.Author.UserID != userId {
		return fmt.Errorf("unauthorized: user does not own this topic")
	}

	existingTopic.Title = req.Title
	existingTopic.Description = req.Description
	existingTopic.UpdatedAt = time.Now()

	if err := s.topicRepo.UpdateTopic(ctx, existingTopic); err != nil {
		log.Printf("[Err] Error updating topic in TopicService.UpdateTopic: %v", err)
		return fmt.Errorf("failed to update topic: %w", err)
	}

	return nil
}

func (s *TopicService) DeleteTopic(ctx context.Context, slug string, userId uint64) error {
	existingTopic, err := s.topicRepo.GetTopicBySlug(ctx, slug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("topic not found")
		}
		log.Printf("[Err] Error getting topic in TopicService.DeleteTopic: %v", err)
		return fmt.Errorf("failed to get topic: %w", err)
	}

	if existingTopic.Author.UserID != userId {
		return fmt.Errorf("unauthorized: user does not own this topic")
	}

	// Check if topic has lessons
	count, err := s.lessonRepo.CountLessonsByTopicID(ctx, existingTopic.ID.Hex())
	if err != nil {
		log.Printf("[Err] Error counting lessons in TopicService.DeleteTopic: %v", err)
		return fmt.Errorf("failed to check lessons: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete topic: has lessons")
	}

	if err := s.topicRepo.DeleteTopic(ctx, slug); err != nil {
		log.Printf("[Err] Error deleting topic in TopicService.DeleteTopic: %v", err)
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	return nil
}

func (s *TopicService) GetLessonsInTopic(ctx context.Context, slug string, page, limit int) ([]*response.LessonResponse, int64, error) {
	topic, err := s.topicRepo.GetTopicBySlug(ctx, slug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, 0, fmt.Errorf("topic not found")
		}
		log.Printf("[Err] Error getting topic in TopicService.GetLessonsInTopic: %v", err)
		return nil, 0, fmt.Errorf("failed to get topic: %w", err)
	}

	lessons, total, err := s.lessonRepo.GetLessonsByTopicID(ctx, topic.ID.Hex(), page, limit)
	if err != nil {
		log.Printf("[Err] Error getting lessons in TopicService.GetLessonsInTopic: %v", err)
		return nil, 0, fmt.Errorf("failed to get lessons: %w", err)
	}

	lessonResponses := response.NewLessonListResponse(lessons)
	return lessonResponses, total, nil
}
