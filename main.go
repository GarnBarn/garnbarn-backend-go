package main

import (
	"github.com/GarnBarn/garnbarn-backend-go/config"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	appConfig config.Config
)

func init() {
	appConfig = config.Load()
}

func main() {
	logger.InitLogger(logger.Config{
		Env: appConfig.Env,
	})

	logrus.Info("Hello")
}
