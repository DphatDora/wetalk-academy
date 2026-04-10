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

type ContentHandler struct {
	contentService *service.ContentService
}

func NewContentHandler(contentService *service.ContentService) *ContentHandler {
	return &ContentHandler{contentService: contentService}
}

// CreateContent
// @Summary Create content for a lesson
// @Description Create lesson content by lesson slug. Requires authentication and lesson ownership.
// @Tags contents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Lesson slug"
// @Param request body request.CreateContentRequest true "Create content payload"
// @Success 201 {object} response.APIResponse "Content created successfully"
// @Failure 400 {object} response.APIResponse "Invalid request format"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden - no permission to add content"
// @Failure 404 {object} response.APIResponse "Lesson not found"
// @Failure 409 {object} response.APIResponse "Content already exists for this lesson"
// @Failure 500 {object} response.APIResponse "Failed to create content"
// @Router /api/v1/lessons/{slug}/content [post]
func (h *ContentHandler) CreateContent(c *gin.Context) {
	ctx := c.Request.Context()
	lessonSlug := c.Param("slug")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in ContentHandler.CreateContent", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in ContentHandler.CreateContent: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	err = h.contentService.CreateContent(ctx, lessonSlug, userID, &req)
	if err != nil {
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
				Message: "You don't have permission to add content to this lesson",
			})
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, response.APIResponse{
				Success: false,
				Message: "Content already exists for this lesson",
			})
			return
		}
		if strings.Contains(err.Error(), "section") || strings.Contains(err.Error(), "required") {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer ContentHandler.CreateContent: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create content",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Content created successfully")
	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Content created successfully",
	})
}

// GetContent
// @Summary Get content by lesson slug
// @Description Retrieve lesson content using lesson slug.
// @Tags contents
// @Accept json
// @Produce json
// @Param slug path string true "Lesson slug"
// @Success 200 {object} response.APIResponse "Content retrieved successfully"
// @Failure 404 {object} response.APIResponse "Lesson or content not found"
// @Failure 500 {object} response.APIResponse "Failed to get content"
// @Router /api/v1/lessons/{slug}/content [get]
func (h *ContentHandler) GetContent(c *gin.Context) {
	ctx := c.Request.Context()
	lessonSlug := c.Param("slug")

	content, err := h.contentService.GetContentByLessonSlug(ctx, lessonSlug)
	if err != nil {
		if strings.Contains(err.Error(), "lesson not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Lesson not found",
			})
			return
		}
		if strings.Contains(err.Error(), "content not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Content not found",
			})
			return
		}

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer ContentHandler.GetContent: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get content",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Content retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Content retrieved successfully",
		Data:    content,
	})
}

// UpdateContent
// @Summary Update lesson content
// @Description Update content of a lesson by lesson slug. Requires authentication and lesson ownership.
// @Tags contents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Lesson slug"
// @Param request body request.UpdateContentRequest true "Update content payload"
// @Success 200 {object} response.APIResponse "Content updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid request format"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden - no permission to update content"
// @Failure 404 {object} response.APIResponse "Lesson or content not found"
// @Failure 500 {object} response.APIResponse "Failed to update content"
// @Router /api/v1/lessons/{slug}/content [put]
func (h *ContentHandler) UpdateContent(c *gin.Context) {
	ctx := c.Request.Context()
	lessonSlug := c.Param("slug")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in ContentHandler.UpdateContent", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in ContentHandler.UpdateContent: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	err = h.contentService.UpdateContent(ctx, lessonSlug, userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "lesson not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Lesson not found",
			})
			return
		}
		if strings.Contains(err.Error(), "content not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Content not found",
			})
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to update this content",
			})
			return
		}
		if strings.Contains(err.Error(), "section") || strings.Contains(err.Error(), "required") {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer ContentHandler.UpdateContent: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update content",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Content updated successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Content updated successfully",
	})
}

// DeleteContent
// @Summary Delete lesson content
// @Description Delete content of a lesson by lesson slug. Requires authentication and lesson ownership.
// @Tags contents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Lesson slug"
// @Success 200 {object} response.APIResponse "Content deleted successfully"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden - no permission to delete content"
// @Failure 404 {object} response.APIResponse "Lesson or content not found"
// @Failure 500 {object} response.APIResponse "Failed to delete content"
// @Router /api/v1/lessons/{slug}/content [delete]
func (h *ContentHandler) DeleteContent(c *gin.Context) {
	ctx := c.Request.Context()
	lessonSlug := c.Param("slug")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in ContentHandler.DeleteContent", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	err = h.contentService.DeleteContent(ctx, lessonSlug, userID)
	if err != nil {
		if strings.Contains(err.Error(), "lesson not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Lesson not found",
			})
			return
		}
		if strings.Contains(err.Error(), "content not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Content not found",
			})
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to delete this content",
			})
			return
		}

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer ContentHandler.DeleteContent: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete content",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Content deleted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Content deleted successfully",
	})
}
