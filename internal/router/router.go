package router

import (
	"felix.chen/login/internal/logger"
	"felix.chen/login/internal/middleware"
	"github.com/gin-gonic/gin"
)

var log = logger.GetLogger()

func CreateRouter() *gin.Engine {
	r := gin.New()
	r.UseRawPath = true // prevent path from being automatically decoded
	r.Use(middleware.Logger(log), gin.Recovery())

	r.GET("/healthcheck", healthcheck)

	return r
}
