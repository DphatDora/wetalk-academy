package service

import (
	"context"
	"fmt"
	"time"
	"wetalk-academy/internal/domain/model"
	"wetalk-academy/internal/domain/repository"
	"wetalk-academy/internal/interface/dto/request"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/package/logger"
	"wetalk-academy/package/util"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LessonService struct {
	lessonRepo  repository.LessonRepository
	topicRepo   repository.TopicRepository
	contentRepo repository.ContentRepository
}

func NewLessonService(lessonRepo repository.LessonRepository, topicRepo repository.TopicRepository, contentRepo repository.ContentRepository) *LessonService {
	return &LessonService{
		lessonRepo:  lessonRepo,
		topicRepo:   topicRepo,
		contentRepo: contentRepo,
	}
}

func (s *LessonService) CreateLesson(ctx context.Context, userId uint64, req *request.CreateLessonRequest) error {
	topic, err := s.topicRepo.GetTopicBySlug(ctx, req.TopicSlug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("topic not found")
		}
		logger.ErrorfWithCtx(ctx, "[Err] Error getting topic in LessonService.CreateLesson: %v", err)
		return fmt.Errorf("failed to get topic: %w", err)
	}

	if topic.Author.UserID != userId {
		return fmt.Errorf("unauthorized: user does not own this topic")
	}

	now := time.Now()
	lesson := &model.Lesson{
		TopicID:    topic.ID,
		Slug:       util.GenerateSlug(req.Title),
		Title:      req.Title,
		OrderIndex: req.OrderIndex,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.lessonRepo.CreateLesson(ctx, lesson); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error creating lesson in LessonService.CreateLesson: %v", err)
		return fmt.Errorf("failed to create lesson: %w", err)
	}

	return nil
}

func (s *LessonService) GetLessonBySlug(ctx context.Context, slug string) (*response.LessonResponse, error) {
	lesson, err := s.lessonRepo.GetLessonBySlug(ctx, slug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("lesson not found")
		}
		logger.ErrorfWithCtx(ctx, "[Err] Error getting lesson in LessonService.GetLessonBySlug: %v", err)
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	return response.NewLessonResponse(lesson), nil
}

func (s *LessonService) UpdateLesson(ctx context.Context, slug string, userId uint64, req *request.UpdateLessonRequest) error {
	existingLesson, err := s.lessonRepo.GetLessonBySlug(ctx, slug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("lesson not found")
		}
		logger.ErrorfWithCtx(ctx, "[Err] Error getting lesson in LessonService.UpdateLesson: %v", err)
		return fmt.Errorf("failed to get lesson: %w", err)
	}

	topicID := existingLesson.TopicID.Hex()
	topic, err := s.topicRepo.GetTopicByID(ctx, topicID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("topic not found")
		}
		logger.ErrorfWithCtx(ctx, "[Err] Error getting topic in LessonService.UpdateLesson: %v", err)
		return fmt.Errorf("failed to get topic: %w", err)
	}

	if topic.Author.UserID != userId {
		return fmt.Errorf("unauthorized: user does not own this topic")
	}

	existingLesson.Title = req.Title
	existingLesson.OrderIndex = req.OrderIndex
	existingLesson.UpdatedAt = time.Now()

	if err := s.lessonRepo.UpdateLesson(ctx, existingLesson); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error updating lesson in LessonService.UpdateLesson: %v", err)
		return fmt.Errorf("failed to update lesson: %w", err)
	}

	return nil
}

func (s *LessonService) DeleteLesson(ctx context.Context, slug string, userId uint64) error {
	existingLesson, err := s.lessonRepo.GetLessonBySlug(ctx, slug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("lesson not found")
		}
		logger.ErrorfWithCtx(ctx, "[Err] Error getting lesson in LessonService.DeleteLesson: %v", err)
		return fmt.Errorf("failed to get lesson: %w", err)
	}

	topicID := existingLesson.TopicID.Hex()
	topic, err := s.topicRepo.GetTopicByID(ctx, topicID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("topic not found")
		}
		logger.ErrorfWithCtx(ctx, "[Err] Error getting topic in LessonService.DeleteLesson: %v", err)
		return fmt.Errorf("failed to get topic: %w", err)
	}

	if topic.Author.UserID != userId {
		return fmt.Errorf("unauthorized: user does not own this topic")
	}

	lessonID := existingLesson.ID.Hex()
	contentExists, err := s.contentRepo.ContentExistsByLessonID(ctx, lessonID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error checking content existence in LessonService.DeleteLesson: %v", err)
		return fmt.Errorf("failed to check content: %w", err)
	}

	if contentExists {
		if err := s.contentRepo.DeleteContentByLessonID(ctx, lessonID); err != nil {
			logger.ErrorfWithCtx(ctx, "[Err] Error deleting content in LessonService.DeleteLesson: %v", err)
			return fmt.Errorf("failed to delete content: %w", err)
		}
	}

	if err := s.lessonRepo.DeleteLesson(ctx, slug); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error deleting lesson in LessonService.DeleteLesson: %v", err)
		return fmt.Errorf("failed to delete lesson: %w", err)
	}

	return nil
}
