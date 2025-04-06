package api

import (
	"timezone-utils/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/health", HealthCheck)
	r.POST("/working-hours/check", CheckWorkingHours)

	return r
}
