package request

type QuizQuestionRequest struct {
	Question      string   `json:"question" binding:"required"`
	Point         int      `json:"point" binding:"required,min=1"`
	Options       []string `json:"options" binding:"required,min=2"`
	CorrectAnswer string   `json:"correctAnswer" binding:"required"`
}

type CreateQuizRequest struct {
	LessonSlug string               `json:"lessonSlug" binding:"required"`
	Title      string               `json:"title" binding:"required"`
	Questions  []QuizQuestionRequest `json:"questions" binding:"required,min=1"`
	TimeLimit  int                  `json:"timeLimit" binding:"required,min=1"`
}

type UpdateQuizRequest struct {
	Title     string               `json:"title" binding:"required"`
	Questions []QuizQuestionRequest `json:"questions" binding:"required,min=1"`
	TimeLimit int                  `json:"timeLimit" binding:"required,min=1"`
}

type SubmitQuizRequest struct {
	QuizID    string   `json:"quizId" binding:"required"`
	Answers   []string `json:"answers" binding:"required"`
	TotalTime int      `json:"totalTime" binding:"required,min=0"`
}
