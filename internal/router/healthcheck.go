package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Healthcheck returns result
// @Summary Health check
// @Description Check the health status of service
// @Produce  json
// @Success 200 {string} string	"Success massage"
// @Failure 500 {string} string "Error message"
// @Router /healthcheck [get]
func healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
