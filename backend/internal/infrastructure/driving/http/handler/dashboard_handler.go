package handler

import (
	"net/http"
	"zentinel/internal/application/port"
	"zentinel/internal/infrastructure/driving/http/dto/response"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardUseCase port.DashboardUseCase
}

func NewDashboardHandler(dashboardUseCase port.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{
		dashboardUseCase: dashboardUseCase,
	}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.dashboardUseCase.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to get dashboard stats", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(stats))
}
