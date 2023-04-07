package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExampleHandler struct{}

func NewExampleHandler() ExampleHandler {
	return ExampleHandler{}
}

func (e *ExampleHandler) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}
