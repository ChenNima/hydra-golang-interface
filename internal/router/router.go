package router

import (
	"felix.chen/login/internal/logger"
	"felix.chen/login/internal/middleware"
	"felix.chen/login/internal/service"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
)

var log = logger.GetLogger()

func CreateRouter() *gin.Engine {
	r := gin.New()
	r.UseRawPath = true // prevent path from being automatically decoded

	r.LoadHTMLGlob("templates/*")
	r.LoadHTMLGlob("resources/views/gin-gonic/*")
	r.Static("resources", "./resources")

	r.Use(middleware.Logger(log), gin.Recovery())

	r.GET("/healthcheck", healthcheck)
	r.GET("/login", loginPage)
	r.POST("/login", doLogin)

	r.GET("/consent", consentPage)
	r.POST("/consent", doConsent)

	{
		auth := r.Group("/auth")
		ab := service.GetAuthboss()
		auth.Use(adapter.Wrap(ab.LoadClientStateMiddleware))
		// auth.Any("/*w", gin.WrapH(http.StripPrefix("/auth", ab.Config.Core.Router)))
		// auth.Any("/*w", gin.WrapH(ab.Config.Core.Router))
		auth.GET("/login", gin.WrapH(ab.Config.Core.Router))
	}

	return r
}
