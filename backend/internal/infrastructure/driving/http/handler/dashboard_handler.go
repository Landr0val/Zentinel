package handler

import (
	"net/http"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/infrastructure/driving/http/dto/response"
	"zentinel/internal/infrastructure/driving/http/errorhandler"

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

// GetStats godoc
// @Summary      Get dashboard statistics
// @Description  Get aggregated statistics for the dashboard
// @Tags         dashboard
// @Produce      json
// @Success      200  {object}  response.DashboardStatsResponseWrapper
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/dashboard/stats [get]
// @Security     BearerAuth
func (h *DashboardHandler) GetStats(c *gin.Context) {
	var stats *dto.DashboardStats
	var err error

	stats, err = h.dashboardUseCase.GetStats(c.Request.Context())
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(stats))
}
