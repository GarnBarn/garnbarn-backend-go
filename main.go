package main

import (
	"fmt"

	"github.com/GarnBarn/garnbarn-backend-go/config"
	"github.com/GarnBarn/garnbarn-backend-go/handler"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/httpserver"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	appConfig config.Config
)

func init() {
	appConfig = config.Load()
	logger.InitLogger(logger.Config{
		Env: appConfig.Env,
	})

}

func main() {
	// Create the http server
	httpServer := httpserver.NewHttpServer()

	// Init the handler
	exampleHandler := handler.NewExampleHandler()

	// Add Routes
	httpServer.GET("/example", exampleHandler.HelloWorld)

	logrus.Info("Listening and serving HTTP on :", appConfig.HTTP_SERVER_PORT)
	httpServer.Run(fmt.Sprint(":", appConfig.HTTP_SERVER_PORT))
}
