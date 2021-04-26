package main

import "felix.chen/login/internal/server"

func main() {
	srv := server.CreateServer()
	server.GracefulShutdown(srv)
}
