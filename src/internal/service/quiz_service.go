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

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type QuizService struct {
	quizRepo       repository.QuizRepository
	submissionRepo repository.QuizSubmissionRepository
	lessonRepo     repository.LessonRepository
	topicRepo      repository.TopicRepository
}

func NewQuizService(
	quizRepo repository.QuizRepository,
	submissionRepo repository.QuizSubmissionRepository,
	lessonRepo repository.LessonRepository,
	topicRepo repository.TopicRepository,
) *QuizService {
	return &QuizService{
		quizRepo:       quizRepo,
		submissionRepo: submissionRepo,
		lessonRepo:     lessonRepo,
		topicRepo:      topicRepo,
	}
}

func (s *QuizService) CreateQuiz(ctx context.Context, userID uint64, req *request.CreateQuizRequest) error {
	lesson, err := s.lessonRepo.GetLessonBySlug(ctx, req.LessonSlug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("lesson not found")
		}
		logger.Errorf("[Err] Error getting lesson in QuizService.CreateQuiz: %v", err)
		return fmt.Errorf("failed to get lesson: %w", err)
	}

	if err := s.checkTopicOwnership(ctx, lesson.TopicID.Hex(), userID); err != nil {
		return err
	}

	questions := make([]model.QuizQuestion, len(req.Questions))
	for i, q := range req.Questions {
		questions[i] = model.QuizQuestion{
			Question:      q.Question,
			Point:         q.Point,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer,
		}
	}

	quiz := &model.Quiz{
		LessonID:  lesson.ID,
		Title:     req.Title,
		Questions: questions,
		TimeLimit: req.TimeLimit,
		CreatedAt: time.Now(),
	}

	if err := s.quizRepo.CreateQuiz(ctx, quiz); err != nil {
		logger.Errorf("[Err] Error creating quiz in QuizService.CreateQuiz: %v", err)
		return fmt.Errorf("failed to create quiz: %w", err)
	}

	return nil
}

func (s *QuizService) GetQuizByID(ctx context.Context, quizID string) (*response.QuizResponse, error) {
	quiz, err := s.quizRepo.GetQuizByID(ctx, quizID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("quiz not found")
		}
		logger.Errorf("[Err] Error getting quiz in QuizService.GetQuizByID: %v", err)
		return nil, fmt.Errorf("failed to get quiz: %w", err)
	}

	return response.NewQuizResponse(quiz), nil
}

func (s *QuizService) GetQuizzesByLessonSlug(ctx context.Context, lessonSlug string) ([]*response.QuizSummaryResponse, error) {
	lesson, err := s.lessonRepo.GetLessonBySlug(ctx, lessonSlug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("lesson not found")
		}
		logger.Errorf("[Err] Error getting lesson in QuizService.GetQuizzesByLessonSlug: %v", err)
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	quizzes, err := s.quizRepo.GetQuizzesByLessonID(ctx, lesson.ID.Hex())
	if err != nil {
		logger.Errorf("[Err] Error getting quizzes in QuizService.GetQuizzesByLessonSlug: %v", err)
		return nil, fmt.Errorf("failed to get quizzes: %w", err)
	}

	return response.NewQuizSummaryListResponse(quizzes), nil
}

func (s *QuizService) UpdateQuiz(ctx context.Context, quizID string, userID uint64, req *request.UpdateQuizRequest) error {
	quiz, err := s.quizRepo.GetQuizByID(ctx, quizID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("quiz not found")
		}
		logger.Errorf("[Err] Error getting quiz in QuizService.UpdateQuiz: %v", err)
		return fmt.Errorf("failed to get quiz: %w", err)
	}

	lesson, err := s.lessonRepo.GetLessonByID(ctx, quiz.LessonID.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("lesson not found")
		}
		logger.Errorf("[Err] Error getting lesson in QuizService.UpdateQuiz: %v", err)
		return fmt.Errorf("failed to get lesson: %w", err)
	}

	if err := s.checkTopicOwnership(ctx, lesson.TopicID.Hex(), userID); err != nil {
		return err
	}

	quiz.Title = req.Title
	questions := make([]model.QuizQuestion, len(req.Questions))
	for i, q := range req.Questions {
		questions[i] = model.QuizQuestion{
			Question:      q.Question,
			Point:         q.Point,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer,
		}
	}
	quiz.Questions = questions
	quiz.TimeLimit = req.TimeLimit

	if err := s.quizRepo.UpdateQuiz(ctx, quiz); err != nil {
		logger.Errorf("[Err] Error updating quiz in QuizService.UpdateQuiz: %v", err)
		return fmt.Errorf("failed to update quiz: %w", err)
	}

	return nil
}

func (s *QuizService) DeleteQuiz(ctx context.Context, quizID string, userID uint64) error {
	quiz, err := s.quizRepo.GetQuizByID(ctx, quizID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("quiz not found")
		}
		logger.Errorf("[Err] Error getting quiz in QuizService.DeleteQuiz: %v", err)
		return fmt.Errorf("failed to get quiz: %w", err)
	}

	lesson, err := s.lessonRepo.GetLessonByID(ctx, quiz.LessonID.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("lesson not found")
		}
		logger.Errorf("[Err] Error getting lesson in QuizService.DeleteQuiz: %v", err)
		return fmt.Errorf("failed to get lesson: %w", err)
	}

	if err := s.checkTopicOwnership(ctx, lesson.TopicID.Hex(), userID); err != nil {
		return err
	}

	if err := s.quizRepo.DeleteQuiz(ctx, quizID); err != nil {
		logger.Errorf("[Err] Error deleting quiz in QuizService.DeleteQuiz: %v", err)
		return fmt.Errorf("failed to delete quiz: %w", err)
	}

	return nil
}

func (s *QuizService) SubmitQuiz(ctx context.Context, userID uint64, req *request.SubmitQuizRequest) (*response.QuizSubmissionResponse, error) {
	quiz, err := s.quizRepo.GetQuizByID(ctx, req.QuizID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("quiz not found")
		}
		logger.Errorf("[Err] Error getting quiz in QuizService.SubmitQuiz: %v", err)
		return nil, fmt.Errorf("failed to get quiz: %w", err)
	}

	quizObjectID, err := bson.ObjectIDFromHex(req.QuizID)
	if err != nil {
		return nil, fmt.Errorf("invalid quiz ID")
	}

	totalScore := s.calculateScore(quiz, req.Answers)

	submission := &model.QuizSubmission{
		QuizID:      quizObjectID,
		UserID:      userID,
		Answers:     req.Answers,
		TotalTime:   req.TotalTime,
		TotalScore:  totalScore,
		SubmittedAt: time.Now(),
	}

	if err := s.submissionRepo.CreateSubmission(ctx, submission); err != nil {
		logger.Errorf("[Err] Error creating submission in QuizService.SubmitQuiz: %v", err)
		return nil, fmt.Errorf("failed to submit quiz: %w", err)
	}

	return response.NewQuizSubmissionResponse(submission), nil
}

func (s *QuizService) checkTopicOwnership(ctx context.Context, topicID string, userID uint64) error {
	topic, err := s.topicRepo.GetTopicByID(ctx, topicID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("topic not found")
		}
		logger.Errorf("[Err] Error getting topic in QuizService.checkTopicOwnership: %v", err)
		return fmt.Errorf("failed to get topic: %w", err)
	}

	if topic.Author.UserID != userID {
		return fmt.Errorf("unauthorized: user does not own this lesson's topic")
	}

	return nil
}

func (s *QuizService) calculateScore(quiz *model.Quiz, answers []string) float64 {
	if len(quiz.Questions) == 0 {
		return 0
	}

	totalPoints := 0
	earnedPoints := 0

	for i, question := range quiz.Questions {
		totalPoints += question.Point
		if i < len(answers) && answers[i] == question.CorrectAnswer {
			earnedPoints += question.Point
		}
	}

	if totalPoints == 0 {
		return 0
	}

	return float64(earnedPoints) / float64(totalPoints) * 100
}
