package httpserver

import (
	"github.com/gin-gonic/gin"
)

func NewHttpServer() *gin.Engine {
	router := gin.Default()
	router.GET("/healthcheck", healthCheckHandler)
	return router
}
