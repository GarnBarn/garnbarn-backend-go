package handler

import (
	"net/http"

	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-gonic/gin"
)

type ExampleHandler struct {
	exampleService service.Example
}

func NewExampleHandler(exampleService service.Example) ExampleHandler {
	return ExampleHandler{
		exampleService: exampleService,
	}
}

func (e *ExampleHandler) HelloWorld(c *gin.Context) {
	err := e.exampleService.Example()

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"hello": "success",
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
	return
}
