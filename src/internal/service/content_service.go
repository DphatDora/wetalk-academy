package service

import (
	"context"
	"fmt"
	"time"
	"wetalk-academy/internal/domain/model"
	"wetalk-academy/internal/domain/repository"
	"wetalk-academy/internal/interface/dto/request"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/package/constant"
	"wetalk-academy/package/logger"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ContentService struct {
	contentRepo repository.ContentRepository
	lessonRepo  repository.LessonRepository
	topicRepo   repository.TopicRepository
}

func NewContentService(
	contentRepo repository.ContentRepository,
	lessonRepo repository.LessonRepository,
	topicRepo repository.TopicRepository,
) *ContentService {
	return &ContentService{
		contentRepo: contentRepo,
		lessonRepo:  lessonRepo,
		topicRepo:   topicRepo,
	}
}

func (s *ContentService) validateOwnership(ctx context.Context, lessonSlug string, userId uint64) (*model.Lesson, error) {
	// Get lesson
	lesson, err := s.lessonRepo.GetLessonBySlug(ctx, lessonSlug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("lesson not found")
		}
		logger.Errorf("[Err] Error getting lesson in ContentService.validateOwnership: %v", err)
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	topicID := lesson.TopicID.Hex()
	topic, err := s.topicRepo.GetTopicByID(ctx, topicID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("topic not found")
		}
		logger.Errorf("[Err] Error getting topic in ContentService.validateOwnership: %v", err)
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	if topic.Author.UserID != userId {
		return nil, fmt.Errorf("unauthorized: user does not own this topic")
	}

	return lesson, nil
}

func (s *ContentService) validateSections(sections []request.ContentSectionRequest) error {
	if len(sections) == 0 {
		return fmt.Errorf("content section is required")
	}

	for i, section := range sections {
		switch section.Type {
		case constant.ContentSectionTypeText:
			if section.Content == "" {
				return fmt.Errorf("section %d: text content is required", i+1)
			}
		case constant.ContentSectionTypeMedia:
			if section.URL == "" {
				return fmt.Errorf("section %d: media URL is required", i+1)
			}
		case constant.ContentSectionTypeCode:
			if section.Content == "" || section.Language == "" {
				return fmt.Errorf("section %d: code content and language are required", i+1)
			}
		default:
			return fmt.Errorf("section %d: invalid section type '%s'", i+1, section.Type)
		}
	}

	return nil
}

func (s *ContentService) CreateContent(ctx context.Context, lessonSlug string, userId uint64, req *request.CreateContentRequest) error {
	// Validate ownership
	lesson, err := s.validateOwnership(ctx, lessonSlug, userId)
	if err != nil {
		return err
	}

	exists, err := s.contentRepo.ContentExistsByLessonID(ctx, lesson.ID.Hex())
	if err != nil {
		logger.Errorf("[Err] Error checking content existence in ContentService.CreateContent: %v", err)
		return fmt.Errorf("failed to check content: %w", err)
	}
	if exists {
		return fmt.Errorf("content already exists for this lesson")
	}

	// Validate sections
	if err := s.validateSections(req.Sections); err != nil {
		return err
	}

	sections := make([]model.ContentSection, len(req.Sections))
	for i, sectionReq := range req.Sections {
		sections[i] = model.ContentSection{
			ID:       uint64(i + 1),
			Type:     sectionReq.Type,
			Content:  sectionReq.Content,
			Language: sectionReq.Language,
			URL:      sectionReq.URL,
		}
	}

	now := time.Now()
	content := &model.Content{
		LessonID:  lesson.ID,
		Sections:  sections,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.contentRepo.CreateContent(ctx, content); err != nil {
		logger.Errorf("[Err] Error creating content in ContentService.CreateContent: %v", err)
		return fmt.Errorf("failed to create content: %w", err)
	}

	return nil
}

func (s *ContentService) GetContentByLessonSlug(ctx context.Context, lessonSlug string) (*response.ContentResponse, error) {
	lesson, err := s.lessonRepo.GetLessonBySlug(ctx, lessonSlug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("lesson not found")
		}
		logger.Errorf("[Err] Error getting lesson in ContentService.GetContentByLessonSlug: %v", err)
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	content, err := s.contentRepo.GetContentByLessonID(ctx, lesson.ID.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("content not found")
		}
		logger.Errorf("[Err] Error getting content in ContentService.GetContentByLessonSlug: %v", err)
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	return response.NewContentResponse(content), nil
}

func (s *ContentService) UpdateContent(ctx context.Context, lessonSlug string, userId uint64, req *request.UpdateContentRequest) error {
	// Validate ownership
	lesson, err := s.validateOwnership(ctx, lessonSlug, userId)
	if err != nil {
		return err
	}

	existingContent, err := s.contentRepo.GetContentByLessonID(ctx, lesson.ID.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("content not found")
		}
		logger.Errorf("[Err] Error getting content in ContentService.UpdateContent: %v", err)
		return fmt.Errorf("failed to get content: %w", err)
	}

	// Validate sections
	if err := s.validateSections(req.Sections); err != nil {
		return err
	}

	// Update sections
	sections := make([]model.ContentSection, len(req.Sections))
	for i, sectionReq := range req.Sections {
		sections[i] = model.ContentSection{
			ID:       uint64(i + 1),
			Type:     sectionReq.Type,
			Content:  sectionReq.Content,
			Language: sectionReq.Language,
			URL:      sectionReq.URL,
		}
	}

	existingContent.Sections = sections
	existingContent.UpdatedAt = time.Now()

	if err := s.contentRepo.UpdateContent(ctx, existingContent); err != nil {
		logger.Errorf("[Err] Error updating content in ContentService.UpdateContent: %v", err)
		return fmt.Errorf("failed to update content: %w", err)
	}

	return nil
}

func (s *ContentService) DeleteContent(ctx context.Context, lessonSlug string, userId uint64) error {
	// Validate ownership
	lesson, err := s.validateOwnership(ctx, lessonSlug, userId)
	if err != nil {
		return err
	}

	exists, err := s.contentRepo.ContentExistsByLessonID(ctx, lesson.ID.Hex())
	if err != nil {
		logger.Errorf("[Err] Error checking content existence in ContentService.DeleteContent: %v", err)
		return fmt.Errorf("failed to check content: %w", err)
	}
	if !exists {
		return fmt.Errorf("content not found")
	}

	if err := s.contentRepo.DeleteContentByLessonID(ctx, lesson.ID.Hex()); err != nil {
		logger.Errorf("[Err] Error deleting content in ContentService.DeleteContent: %v", err)
		return fmt.Errorf("failed to delete content: %w", err)
	}

	return nil
}
