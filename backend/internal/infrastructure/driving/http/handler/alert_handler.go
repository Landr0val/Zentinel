package handler

import (
	"net/http"
	"strconv"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/domain"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/repository"
	"zentinel/internal/infrastructure/driving/http/dto/response"
	"zentinel/internal/infrastructure/driving/http/errorhandler"

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

// Create godoc
// @Summary      Create a new alert
// @Description  Create a new alert manually
// @Tags         alerts
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateAlertRequest true "Create Alert Request"
// @Success      201  {object}  response.AlertResponseWrapper
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/alerts [post]
// @Security     BearerAuth
func (h *AlertHandler) Create(c *gin.Context) {
	var req dto.CreateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleRequestError(c, err)
		return
	}

	alert, err := h.alertUseCase.CreateAlert(c.Request.Context(), req)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(alert))
}

// Get godoc
// @Summary      Get alert by ID
// @Description  Get alert details by ID
// @Tags         alerts
// @Produce      json
// @Param        id   path      string  true  "Alert ID"
// @Success      200  {object}  response.AlertResponseWrapper
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/alerts/{id} [get]
// @Security     BearerAuth
func (h *AlertHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid alert ID", err.Error()))
		return
	}

	alert, err := h.alertUseCase.GetAlert(c.Request.Context(), id)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	if alert == nil {
		errorhandler.HandleDomainError(c, domain.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(alert))
}

// List godoc
// @Summary      List alerts
// @Description  List alerts with filtering and pagination
// @Tags         alerts
// @Produce      json
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Page size"
// @Param        search     query     string  false  "Search term"
// @Param        status     query     string  false  "Status"
// @Param        severity   query     string  false  "Severity"
// @Success      200        {object}  response.AlertListResponseWrapper
// @Failure      500        {object}  response.ErrorResponse
// @Router       /api/v1/alerts [get]
// @Security     BearerAuth
func (h *AlertHandler) List(c *gin.Context) {
	filter := repository.AlertFilter{
		Page:     1,
		PageSize: 10,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filter.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			filter.PageSize = pageSize
		}
	}

	if search := c.Query("search"); search != "" {
		filter.Search = search
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := enums.AlertStatus(statusStr)
		filter.Status = &status
	}

	if severityStr := c.Query("severity"); severityStr != "" {
		severity := enums.AlertSeverity(severityStr)
		filter.Severity = &severity
	}

	alerts, total, err := h.alertUseCase.ListAlerts(c.Request.Context(), filter)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(alerts, filter.Page, filter.PageSize, total))
}

// UpdateStatus godoc
// @Summary      Update alert status
// @Description  Update the status of an alert
// @Tags         alerts
// @Accept       json
// @Produce      json
// @Param        id      path      string                      true  "Alert ID"
// @Param        request body      dto.UpdateAlertStatusRequest true  "Update Alert Status Request"
// @Success      200     {object}  response.AlertResponseWrapper
// @Failure      400     {object}  response.ErrorResponse
// @Failure      404     {object}  response.ErrorResponse
// @Failure      500     {object}  response.ErrorResponse
// @Router       /api/v1/alerts/{id}/status [patch]
// @Security     BearerAuth
func (h *AlertHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid alert ID", err.Error()))
		return
	}

	var req dto.UpdateAlertStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleRequestError(c, err)
		return
	}

	alert, err := h.alertUseCase.UpdateAlertStatus(c.Request.Context(), id, req)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(alert))
}
