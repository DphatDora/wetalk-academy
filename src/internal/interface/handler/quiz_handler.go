package handler

import (
	"net/http"
	"strings"
	"wetalk-academy/internal/interface/dto/request"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/internal/service"
	"wetalk-academy/package/util"
	"wetalk-academy/package/logger"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizService *service.QuizService
}

func NewQuizHandler(quizService *service.QuizService) *QuizHandler {
	return &QuizHandler{quizService: quizService}
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.Errorf("[Err] %s in QuizHandler.CreateQuiz", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("[Err] Error binding JSON in QuizHandler.CreateQuiz: %v", err)
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

		logger.Errorf("[Err] Error in service layer QuizHandler.CreateQuiz: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create quiz",
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Quiz created successfully",
	})
}

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

		logger.Errorf("[Err] Error in service layer QuizHandler.GetQuizByID: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get quiz",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Quiz retrieved successfully",
		Data:    quiz,
	})
}

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

		logger.Errorf("[Err] Error in service layer QuizHandler.GetQuizzesByLessonSlug: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get quizzes",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Quizzes retrieved successfully",
		Data:    quizzes,
	})
}

func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	quizID := c.Param("id")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.Errorf("[Err] %s in QuizHandler.UpdateQuiz", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.UpdateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("[Err] Error binding JSON in QuizHandler.UpdateQuiz: %v", err)
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

		logger.Errorf("[Err] Error in service layer QuizHandler.UpdateQuiz: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update quiz",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Quiz updated successfully",
	})
}

func (h *QuizHandler) DeleteQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	quizID := c.Param("id")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.Errorf("[Err] %s in QuizHandler.DeleteQuiz", err.Error())
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

		logger.Errorf("[Err] Error in service layer QuizHandler.DeleteQuiz: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete quiz",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Quiz deleted successfully",
	})
}

func (h *QuizHandler) SubmitQuiz(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.Errorf("[Err] %s in QuizHandler.SubmitQuiz", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("[Err] Error binding JSON in QuizHandler.SubmitQuiz: %v", err)
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

		logger.Errorf("[Err] Error in service layer QuizHandler.SubmitQuiz: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to submit quiz",
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Quiz submitted successfully",
		Data:    result,
	})
}
