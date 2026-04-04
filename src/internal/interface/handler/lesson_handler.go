package handler

import (
	"net/http"
	"strings"
	"wetalk-academy/internal/interface/dto/request"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/package/logger"
	"wetalk-academy/internal/service"
	"wetalk-academy/package/util"

	"github.com/gin-gonic/gin"
)

type LessonHandler struct {
	lessonService *service.LessonService
}

func NewLessonHandler(lessonService *service.LessonService) *LessonHandler {
	return &LessonHandler{lessonService: lessonService}
}

func (h *LessonHandler) CreateLesson(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.Errorf("[Err] %s in LessonHandler.CreateLesson", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("[Err] Error binding JSON in LessonHandler.CreateLesson: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	err = h.lessonService.CreateLesson(ctx, userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "topic not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Topic not found",
			})
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to add lessons to this topic",
			})
			return
		}

		logger.Errorf("[Err] Error in service layer LessonHandler.CreateLesson: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create lesson",
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Lesson created successfully",
	})
}

func (h *LessonHandler) GetLessonBySlug(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	lesson, err := h.lessonService.GetLessonBySlug(ctx, slug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Lesson not found",
			})
			return
		}

		logger.Errorf("[Err] Error in service layer LessonHandler.GetLessonBySlug: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get lesson",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Lesson retrieved successfully",
		Data:    lesson,
	})
}

func (h *LessonHandler) UpdateLesson(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.Errorf("[Err] %s in LessonHandler.UpdateLesson", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("[Err] Error binding JSON in LessonHandler.UpdateLesson: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	err = h.lessonService.UpdateLesson(ctx, slug, userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Lesson not found",
			})
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to update this lesson",
			})
			return
		}

		logger.Errorf("[Err] Error in service layer LessonHandler.UpdateLesson: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update lesson",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Lesson updated successfully",
	})
}

func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.Errorf("[Err] %s in LessonHandler.DeleteLesson", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	err = h.lessonService.DeleteLesson(ctx, slug, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Lesson not found",
			})
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to delete this lesson",
			})
			return
		}

		logger.Errorf("[Err] Error in service layer LessonHandler.DeleteLesson: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete lesson",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Lesson deleted successfully",
	})
}
