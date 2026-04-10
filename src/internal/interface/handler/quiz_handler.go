package handler

import (
	"net/http"
	"strings"
	"wetalk-academy/internal/interface/dto/request"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/internal/service"
	"wetalk-academy/package/logger"
	"wetalk-academy/package/util"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizService *service.QuizService
}

func NewQuizHandler(quizService *service.QuizService) *QuizHandler {
	return &QuizHandler{quizService: quizService}
}

// CreateQuiz
// @Summary Create a quiz
// @Description Create a new quiz for a lesson. Requires authentication and lesson ownership.
// @Tags quizzes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateQuizRequest true "Create quiz payload"
// @Success 201 {object} response.APIResponse "Quiz created successfully"
// @Failure 400 {object} response.APIResponse "Invalid request format"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden - no permission to create quiz"
// @Failure 404 {object} response.APIResponse "Lesson not found"
// @Failure 500 {object} response.APIResponse "Failed to create quiz"
// @Router /api/v1/quizzes [post]
func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in QuizHandler.CreateQuiz", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in QuizHandler.CreateQuiz: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	if err := h.quizService.CreateQuiz(ctx, userID, &req); err != nil {
		if strings.Contains(err.Error(), "lesson not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Lesson not found",
			})
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to create a quiz for this lesson",
			})
			return
		}

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer QuizHandler.CreateQuiz: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create quiz",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Quiz created successfully")
	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Quiz created successfully",
	})
}

// GetQuizByID
// @Summary Get quiz by ID
// @Description Retrieve quiz detail by quiz ID.
// @Tags quizzes
// @Accept json
// @Produce json
// @Param id path string true "Quiz ID"
// @Success 200 {object} response.APIResponse "Quiz retrieved successfully"
// @Failure 404 {object} response.APIResponse "Quiz not found"
// @Failure 500 {object} response.APIResponse "Failed to get quiz"
// @Router /api/v1/quizzes/{id} [get]
func (h *QuizHandler) GetQuizByID(c *gin.Context) {
	ctx := c.Request.Context()
	quizID := c.Param("id")

	quiz, err := h.quizService.GetQuizByID(ctx, quizID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Quiz not found",
			})
			return
		}

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer QuizHandler.GetQuizByID: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get quiz",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Quiz retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Quiz retrieved successfully",
		Data:    quiz,
	})
}

// GetQuizzesByLessonSlug
// @Summary Get quizzes by lesson slug
// @Description Retrieve all quizzes of a lesson by lesson slug.
// @Tags quizzes
// @Accept json
// @Produce json
// @Param slug path string true "Lesson slug"
// @Success 200 {object} response.APIResponse "Quizzes retrieved successfully"
// @Failure 404 {object} response.APIResponse "Lesson not found"
// @Failure 500 {object} response.APIResponse "Failed to get quizzes"
// @Router /api/v1/lessons/{slug}/quiz [get]
func (h *QuizHandler) GetQuizzesByLessonSlug(c *gin.Context) {
	ctx := c.Request.Context()
	lessonSlug := c.Param("slug")

	quizzes, err := h.quizService.GetQuizzesByLessonSlug(ctx, lessonSlug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Lesson not found",
			})
			return
		}

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer QuizHandler.GetQuizzesByLessonSlug: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get quizzes",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Quizzes retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Quizzes retrieved successfully",
		Data:    quizzes,
	})
}

// UpdateQuiz
// @Summary Update a quiz
// @Description Update quiz detail by quiz ID. Requires authentication and quiz ownership.
// @Tags quizzes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Quiz ID"
// @Param request body request.UpdateQuizRequest true "Update quiz payload"
// @Success 200 {object} response.APIResponse "Quiz updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid request format"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden - no permission to update quiz"
// @Failure 404 {object} response.APIResponse "Quiz not found"
// @Failure 500 {object} response.APIResponse "Failed to update quiz"
// @Router /api/v1/quizzes/{id} [put]
func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	quizID := c.Param("id")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in QuizHandler.UpdateQuiz", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.UpdateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in QuizHandler.UpdateQuiz: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	if err := h.quizService.UpdateQuiz(ctx, quizID, userID, &req); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Quiz not found",
			})
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to update this quiz",
			})
			return
		}

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer QuizHandler.UpdateQuiz: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update quiz",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Quiz updated successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Quiz updated successfully",
	})
}

// DeleteQuiz
// @Summary Delete a quiz
// @Description Delete quiz by quiz ID. Requires authentication and quiz ownership.
// @Tags quizzes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Quiz ID"
// @Success 200 {object} response.APIResponse "Quiz deleted successfully"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden - no permission to delete quiz"
// @Failure 404 {object} response.APIResponse "Quiz not found"
// @Failure 500 {object} response.APIResponse "Failed to delete quiz"
// @Router /api/v1/quizzes/{id} [delete]
func (h *QuizHandler) DeleteQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	quizID := c.Param("id")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in QuizHandler.DeleteQuiz", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	if err := h.quizService.DeleteQuiz(ctx, quizID, userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Quiz not found",
			})
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to delete this quiz",
			})
			return
		}

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer QuizHandler.DeleteQuiz: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete quiz",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Quiz deleted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Quiz deleted successfully",
	})
}

// SubmitQuiz
// @Summary Submit quiz answers
// @Description Submit answers for a quiz and get the submission result. Requires authentication.
// @Tags quizzes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.SubmitQuizRequest true "Submit quiz payload"
// @Success 201 {object} response.APIResponse "Quiz submitted successfully"
// @Failure 400 {object} response.APIResponse "Invalid request format or invalid quiz ID"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 404 {object} response.APIResponse "Quiz not found"
// @Failure 500 {object} response.APIResponse "Failed to submit quiz"
// @Router /api/v1/quizzes/submit [post]
func (h *QuizHandler) SubmitQuiz(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in QuizHandler.SubmitQuiz", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in QuizHandler.SubmitQuiz: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	result, err := h.quizService.SubmitQuiz(ctx, userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "quiz not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Quiz not found",
			})
			return
		}
		if strings.Contains(err.Error(), "invalid quiz ID") {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Invalid quiz ID",
			})
			return
		}

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer QuizHandler.SubmitQuiz: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to submit quiz",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Quiz submitted successfully")
	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Quiz submitted successfully",
		Data:    result,
	})
}
