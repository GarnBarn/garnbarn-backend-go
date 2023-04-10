package handler

import (
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-gonic/gin"
)

type AssignmentHandler struct {
	assignmentService service.AssignmentService
}

func NewAssignmentHandler(rg *gin.RouterGroup, assignmentService service.AssignmentService) AssignmentHandler {
	result := AssignmentHandler{
		assignmentService: assignmentService,
	}

	// router := rg.Group("/assignment")

	// router.POST("/", assignmentService.CreateAssignment)

	return result
}
