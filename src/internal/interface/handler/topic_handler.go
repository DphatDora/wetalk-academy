package handler

import (
	"net/http"
	"strconv"
	"strings"
	"wetalk-academy/internal/interface/dto/request"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/internal/service"
	"wetalk-academy/package/constant"
	"wetalk-academy/package/util"
	"wetalk-academy/package/logger"

	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	topicService *service.TopicService
}

func NewTopicHandler(topicService *service.TopicService) *TopicHandler {
	return &TopicHandler{topicService: topicService}
}

func (h *TopicHandler) CreateTopic(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.Errorf("[Err] %s in TopicHandler.CreateTopic", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("[Err] Error binding JSON in TopicHandler.CreateTopic: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	err = h.topicService.CreateTopic(ctx, userID, &req)
	if err != nil {
		logger.Errorf("[Err] Error in service layer TopicHandler.CreateTopic: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create topic",
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Topic created successfully",
	})
}

func (h *TopicHandler) GetTopics(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	topics, total, err := h.topicService.GetTopics(ctx, page, limit)
	if err != nil {
		logger.Errorf("[Err] Error in service layer TopicHandler.GetTopics: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get topics",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Topics retrieved successfully",
		Data:    topics,
		Pagination: &response.Pagination{
			Total:   total,
			Page:    page,
			Limit:   limit,
			NextURL: util.BuildNextURL(c, total, page, limit),
		},
	})
}

func (h *TopicHandler) GetTopicBySlug(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	topic, err := h.topicService.GetTopicBySlug(ctx, slug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Topic not found",
			})
			return
		}

		logger.Errorf("[Err] Error in service layer TopicHandler.GetTopicBySlug: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get topic",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Topic retrieved successfully",
		Data:    topic,
	})
}

func (h *TopicHandler) UpdateTopic(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.Errorf("[Err] %s in TopicHandler.UpdateTopic", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.UpdateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("[Err] Error binding JSON in TopicHandler.UpdateTopic: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	err = h.topicService.UpdateTopic(ctx, slug, userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Topic not found",
			})
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to update this topic",
			})
			return
		}

		logger.Errorf("[Err] Error in service layer TopicHandler.UpdateTopic: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update topic",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Topic updated successfully",
	})
}

func (h *TopicHandler) DeleteTopic(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.Errorf("[Err] %s in TopicHandler.DeleteTopic", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	err = h.topicService.DeleteTopic(ctx, slug, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Topic not found",
			})
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to delete this topic",
			})
			return
		}
		if strings.Contains(err.Error(), "has lessons") {
			c.JSON(http.StatusConflict, response.APIResponse{
				Success: false,
				Message: "Cannot delete topic with existing lessons",
			})
			return
		}

		logger.Errorf("[Err] Error in service layer TopicHandler.DeleteTopic: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete topic",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Topic deleted successfully",
	})
}

func (h *TopicHandler) GetLessonsInTopic(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	lessons, total, err := h.topicService.GetLessonsInTopic(ctx, slug, page, limit)
	if err != nil {
		if strings.Contains(err.Error(), "topic not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Topic not found",
			})
			return
		}

		logger.Errorf("[Err] Error in service layer TopicHandler.GetLessonsInTopic: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get lessons",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Lessons retrieved successfully",
		Data:    lessons,
		Pagination: &response.Pagination{
			Total:   total,
			Page:    page,
			Limit:   limit,
			NextURL: util.BuildNextURL(c, total, page, limit),
		},
	})
}
