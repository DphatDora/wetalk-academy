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
