package handler

import (
	"net/http"
	"strconv"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/infrastructure/driving/http/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AlertHandler struct {
	alertUseCase port.AlertUseCase
}

func NewAlertHandler(alertUseCase port.AlertUseCase) *AlertHandler {
	return &AlertHandler{
		alertUseCase: alertUseCase,
	}
}

func (h *AlertHandler) Create(c *gin.Context) {
	var req dto.CreateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_REQUEST", "Invalid request body", err.Error()))
		return
	}

	alert, err := h.alertUseCase.CreateAlert(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to create alert", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(alert))
}

func (h *AlertHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid alert ID", err.Error()))
		return
	}

	alert, err := h.alertUseCase.GetAlert(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to get alert", err.Error()))
		return
	}

	if alert == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("NOT_FOUND", "Alert not found", ""))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(alert))
}

func (h *AlertHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	alerts, total, err := h.alertUseCase.ListAlerts(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to list alerts", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(alerts, page, pageSize, total))
}

func (h *AlertHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid alert ID", err.Error()))
		return
	}

	var req dto.UpdateAlertStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_REQUEST", "Invalid request body", err.Error()))
		return
	}

	alert, err := h.alertUseCase.UpdateAlertStatus(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to update alert status", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(alert))
}
