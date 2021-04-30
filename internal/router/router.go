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

	r.LoadHTMLGlob("templates/*")

	r.Use(middleware.Logger(log), gin.Recovery())

	r.GET("/healthcheck", healthcheck)
	r.GET("/login", loginPage)
	r.POST("/login", doLogin)

	r.GET("/consent", consentPage)
	r.POST("/consent", doConsent)

	r.GET("/register", registerPage)
	r.POST("/register", doRegister)

	return r
}
