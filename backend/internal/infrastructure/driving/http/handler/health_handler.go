package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db        *gorm.DB
	startTime time.Time
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now(),
	}
}

// Check godoc
// @Summary      Health Check
// @Description  Checks the health of the service and its dependencies
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	dbStatus := "up"

	sqlDB, err := h.db.DB()
	if err != nil {
		dbStatus = "down"
	} else if err := sqlDB.Ping(); err != nil {
		dbStatus = "down"
	}

	uptime := time.Since(h.startTime).Seconds()

	c.JSON(http.StatusOK, gin.H{
		"status": "up",
		"checks": gin.H{
			"database":       dbStatus,
			"uptime_seconds": int64(uptime),
		},
	})
}
