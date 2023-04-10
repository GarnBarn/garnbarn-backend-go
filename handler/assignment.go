package handler

import (
	"net/http"

	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-gonic/gin"
)

type AssignmentHandler struct {
	assignmentService service.AssignmentService
}

func NewAssignmentHandler(assignmentService service.AssignmentService) AssignmentHandler {
	return AssignmentHandler{
		assignmentService: assignmentService,
	}
}

func (a *AssignmentHandler) AssignmentRoute(rg *gin.RouterGroup) {
	router := rg.Group("/assignment")

	router.POST("/", a.CreateAssignment)
}

func (a *AssignmentHandler) CreateAssignment(c *gin.Context) {
	var assignment model.Assignment

	// Check the conditional operator (?:) later.
	// add validate
	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if result := a.assignmentService.CreateAssignment(&assignment); result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result,
		})
		return
	}

	c.JSON(http.StatusCreated, &assignment)

}
