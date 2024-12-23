// Package main ...
package main

import (
	_ "github.com/KozlovNikolai/pfp/docs"
	"github.com/KozlovNikolai/pfp/internal/chat/transport/httpserver"
	"github.com/KozlovNikolai/pfp/internal/pkg/config"
)

// @title 	Chat Service API
// @version	1.0
// @description Chat service API in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host 	localhost:8443
// @BasePath /
func main() {
	config.MustLoad()

	server := httpserver.NewRouter()

	server.Run()
}
