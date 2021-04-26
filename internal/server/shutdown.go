package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Warn(("Shutdown Server ..."))
}
