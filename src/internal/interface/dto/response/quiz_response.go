package response

import (
	"time"
	"wetalk-academy/internal/domain/model"
)

type QuizQuestionResponse struct {
	Question      string   `json:"question"`
	Point         int      `json:"point"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correctAnswer"`
}

type QuizResponse struct {
	ID        string                 `json:"id"`
	LessonID  string                 `json:"lessonId"`
	Title     string                 `json:"title"`
	Questions []QuizQuestionResponse `json:"questions"`
	TimeLimit int                    `json:"timeLimit"`
	CreatedAt time.Time              `json:"createdAt"`
}

func NewQuizResponse(quiz *model.Quiz) *QuizResponse {
	questions := make([]QuizQuestionResponse, len(quiz.Questions))
	for i, q := range quiz.Questions {
		questions[i] = QuizQuestionResponse{
			Question:      q.Question,
			Point:         q.Point,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer,
		}
	}

	return &QuizResponse{
		ID:        quiz.ID.Hex(),
		LessonID:  quiz.LessonID.Hex(),
		Title:     quiz.Title,
		Questions: questions,
		TimeLimit: quiz.TimeLimit,
		CreatedAt: quiz.CreatedAt,
	}
}

type QuizSummaryResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	TimeLimit int    `json:"timeLimit"`
}

func NewQuizSummaryResponse(quiz *model.Quiz) *QuizSummaryResponse {
	return &QuizSummaryResponse{
		ID:        quiz.ID.Hex(),
		Title:     quiz.Title,
		TimeLimit: quiz.TimeLimit,
	}
}

func NewQuizSummaryListResponse(quizzes []*model.Quiz) []*QuizSummaryResponse {
	result := make([]*QuizSummaryResponse, len(quizzes))
	for i, q := range quizzes {
		result[i] = NewQuizSummaryResponse(q)
	}
	return result
}

type QuizSubmissionResponse struct {
	QuizID      string    `json:"quizId"`
	UserID      uint64    `json:"userId"`
	Answers     []string  `json:"answers"`
	TotalTime   int       `json:"totalTime"`
	TotalScore  float64   `json:"totalScore"`
	SubmittedAt time.Time `json:"submittedAt"`
}

func NewQuizSubmissionResponse(submission *model.QuizSubmission) *QuizSubmissionResponse {
	return &QuizSubmissionResponse{
		QuizID:      submission.QuizID.Hex(),
		UserID:      submission.UserID,
		Answers:     submission.Answers,
		TotalTime:   submission.TotalTime,
		TotalScore:  submission.TotalScore,
		SubmittedAt: submission.SubmittedAt,
	}
}
