package handler

import (
	"log"
	"net/http"
	"wetalk-academy/internal/interface/dto/request"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/internal/service"

	"github.com/gin-gonic/gin"
)

type Judge0Handler struct {
	judge0Service *service.Judge0Service
}

func NewJudge0Handler(judge0Service *service.Judge0Service) *Judge0Handler {
	return &Judge0Handler{judge0Service: judge0Service}
}

func (h *Judge0Handler) SubmitCode(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.SubmitCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in Judge0Handler.SubmitCode: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	result, err := h.judge0Service.SubmitCode(ctx, &req)
	if err != nil {
		log.Printf("[Err] Error in service layer Judge0Handler.SubmitCode: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to submit code",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Code executed successfully",
		Data:    result,
	})
}
