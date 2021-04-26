package server

import (
	"net/http"

	"felix.chen/login/internal/logger"
	"felix.chen/login/internal/router"
	"felix.chen/login/internal/util"
)

var log = logger.GetLogger()

// CreateServer creat gin webserver and setup graceful shutdown
func CreateServer() *http.Server {
	r := router.CreateRouter()
	port := util.GetEnvWithDefault("PORT", "8083")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return srv
}
