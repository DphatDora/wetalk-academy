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

type LessonHandler struct {
	lessonService *service.LessonService
}

func NewLessonHandler(lessonService *service.LessonService) *LessonHandler {
	return &LessonHandler{lessonService: lessonService}
}

// CreateLesson
// @Summary Create a new lesson
// @Description Create a new lesson under a topic. Requires authentication.
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateLessonRequest true "Create lesson payload"
// @Success 201 {object} response.APIResponse "Lesson created successfully"
// @Failure 400 {object} response.APIResponse "Invalid request format"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden - no permission to add lessons to this topic"
// @Failure 404 {object} response.APIResponse "Topic not found"
// @Failure 500 {object} response.APIResponse "Failed to create lesson"
// @Router /api/v1/lessons [post]
func (h *LessonHandler) CreateLesson(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in LessonHandler.CreateLesson", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in LessonHandler.CreateLesson: %v", err)
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

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer LessonHandler.CreateLesson: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create lesson",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Lesson created successfully")
	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Lesson created successfully",
	})
}

// GetLessonBySlug
// @Summary Get lesson by slug
// @Description Retrieve a lesson by its slug.
// @Tags lessons
// @Accept json
// @Produce json
// @Param slug path string true "Lesson slug"
// @Success 200 {object} response.APIResponse "Lesson retrieved successfully"
// @Failure 404 {object} response.APIResponse "Lesson not found"
// @Failure 500 {object} response.APIResponse "Failed to get lesson"
// @Router /api/v1/lessons/{slug} [get]
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

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer LessonHandler.GetLessonBySlug: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get lesson",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Lesson retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Lesson retrieved successfully",
		Data:    lesson,
	})
}

// UpdateLesson
// @Summary Update an existing lesson
// @Description Update an existing lesson by its slug. Requires authentication and ownership.
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Lesson slug"
// @Param request body request.UpdateLessonRequest true "Update lesson payload"
// @Success 200 {object} response.APIResponse "Lesson updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid request format"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden - no permission to update this lesson"
// @Failure 404 {object} response.APIResponse "Lesson not found"
// @Failure 500 {object} response.APIResponse "Failed to update lesson"
// @Router /api/v1/lessons/{slug} [put]
func (h *LessonHandler) UpdateLesson(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in LessonHandler.UpdateLesson", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in LessonHandler.UpdateLesson: %v", err)
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

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer LessonHandler.UpdateLesson: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update lesson",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Lesson updated successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Lesson updated successfully",
	})
}

// DeleteLesson
// @Summary Delete a lesson
// @Description Delete a lesson by its slug. Requires authentication and ownership.
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Lesson slug"
// @Success 200 {object} response.APIResponse "Lesson deleted successfully"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden - no permission to delete this lesson"
// @Failure 404 {object} response.APIResponse "Lesson not found"
// @Failure 500 {object} response.APIResponse "Failed to delete lesson"
// @Router /api/v1/lessons/{slug} [delete]
func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in LessonHandler.DeleteLesson", err.Error())
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

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer LessonHandler.DeleteLesson: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete lesson",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Lesson deleted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Lesson deleted successfully",
	})
}
