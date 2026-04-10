package handler

import (
	"net/http"
	"strconv"
	"strings"
	"wetalk-academy/internal/interface/dto/request"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/internal/service"
	"wetalk-academy/package/constant"
	"wetalk-academy/package/logger"
	"wetalk-academy/package/util"

	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	topicService *service.TopicService
}

func NewTopicHandler(topicService *service.TopicService) *TopicHandler {
	return &TopicHandler{topicService: topicService}
}

// CreateTopic
// @Summary Create a new topic
// @Description Create a new topic with title, description and author information. Requires authentication.
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateTopicRequest true "Create topic payload"
// @Success 201 {object} response.APIResponse "Topic created successfully"
// @Failure 400 {object} response.APIResponse "Invalid request format"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 500 {object} response.APIResponse "Failed to create topic"
// @Router /api/v1/topics [post]
func (h *TopicHandler) CreateTopic(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in TopicHandler.CreateTopic", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in TopicHandler.CreateTopic: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	err = h.topicService.CreateTopic(ctx, userID, &req)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer TopicHandler.CreateTopic: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create topic",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Topic created successfully")
	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Topic created successfully",
	})
}

// GetTopics
// @Summary Get list of topics
// @Description Retrieve a paginated list of topics.
// @Tags topics
// @Accept json
// @Produce json
// @Param page query int false "Page number" minimum(1)
// @Param limit query int false "Items per page" minimum(1) maximum(100)
// @Success 200 {object} response.APIResponse{data=[]response.TopicResponse} "Topics retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to get topics"
// @Router /api/v1/topics [get]
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
		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer TopicHandler.GetTopics: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get topics",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Topics retrieved successfully")
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

// GetTopicBySlug
// @Summary Get topic by slug
// @Description Retrieve a single topic by its slug.
// @Tags topics
// @Accept json
// @Produce json
// @Param slug path string true "Topic slug"
// @Success 200 {object} response.APIResponse{data=response.TopicResponse} "Topic retrieved successfully"
// @Failure 404 {object} response.APIResponse "Topic not found"
// @Failure 500 {object} response.APIResponse "Failed to get topic"
// @Router /api/v1/topics/{slug} [get]
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

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer TopicHandler.GetTopicBySlug: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get topic",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Topic retrieved successfully by slug")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Topic retrieved successfully",
		Data:    topic,
	})
}

// UpdateTopic
// @Summary Update a topic
// @Description Update an existing topic by slug. Requires authentication and ownership of the topic.
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Topic slug"
// @Param request body request.UpdateTopicRequest true "Update topic payload"
// @Success 200 {object} response.APIResponse "Topic updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid request format"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden - no permission to update this topic"
// @Failure 404 {object} response.APIResponse "Topic not found"
// @Failure 500 {object} response.APIResponse "Failed to update topic"
// @Router /api/v1/topics/{slug} [put]
func (h *TopicHandler) UpdateTopic(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in TopicHandler.UpdateTopic", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.UpdateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in TopicHandler.UpdateTopic: %v", err)
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

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer TopicHandler.UpdateTopic: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update topic",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Topic updated successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Topic updated successfully",
	})
}

// DeleteTopic
// @Summary Delete a topic
// @Description Delete an existing topic by slug. Requires authentication and ownership of the topic. Topic cannot be deleted if it has existing lessons.
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Topic slug"
// @Success 200 {object} response.APIResponse "Topic deleted successfully"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden - no permission to delete this topic"
// @Failure 404 {object} response.APIResponse "Topic not found"
// @Failure 409 {object} response.APIResponse "Cannot delete topic with existing lessons"
// @Failure 500 {object} response.APIResponse "Failed to delete topic"
// @Router /api/v1/topics/{slug} [delete]
func (h *TopicHandler) DeleteTopic(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in TopicHandler.DeleteTopic", err.Error())
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

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer TopicHandler.DeleteTopic: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete topic",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Topic deleted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Topic deleted successfully",
	})
}

// GetLessonsInTopic
// @Summary Get lessons in topic
// @Description Retrieve a paginated list of lessons belonging to a topic identified by slug.
// @Tags topics
// @Accept json
// @Produce json
// @Param slug path string true "Topic slug"
// @Param page query int false "Page number" minimum(1)
// @Param limit query int false "Items per page" minimum(1) maximum(100)
// @Success 200 {object} response.APIResponse "Lessons retrieved successfully"
// @Failure 404 {object} response.APIResponse "Topic not found"
// @Failure 500 {object} response.APIResponse "Failed to get lessons"
// @Router /api/v1/topics/{slug}/lessons [get]
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

		logger.ErrorfWithCtx(ctx, "[Err] Error in service layer TopicHandler.GetLessonsInTopic: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get lessons",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Lessons in topic retrieved successfully")
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
