package handler

import (
	"net/http"
	"wetalk-academy/internal/interface/dto/request"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/internal/service"
	"wetalk-academy/package/logger"

	"github.com/gin-gonic/gin"
)

type Judge0Handler struct {
	judge0Service *service.Judge0Service
}

func NewJudge0Handler(judge0Service *service.Judge0Service) *Judge0Handler {
	return &Judge0Handler{judge0Service: judge0Service}
}

// SubmitCode
// @Summary Submit code for execution
// @Description Submit code to be executed by Judge0 API.
// @Tags judge0
// @Accept json
// @Produce json
// @Param request body request.SubmitCodeRequest true "Submit code payload"
// @Success 200 {object} response.APIResponse "Code executed successfully"
// @Failure 400 {object} response.APIResponse "Invalid request format"
// @Failure 500 {object} response.APIResponse "Failed to submit code"
// @Router /api/v1/judge0/submit [post]
func (h *Judge0Handler) SubmitCode(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.SubmitCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("[Err] Error binding JSON in Judge0Handler.SubmitCode: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	result, err := h.judge0Service.SubmitCode(ctx, &req)
	if err != nil {
		logger.Errorf("[Err] Error in service layer Judge0Handler.SubmitCode: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to submit code",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Code executed successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Code executed successfully",
		Data:    result,
	})
}
