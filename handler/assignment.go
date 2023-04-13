package handler

import (
	"net/http"

	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AssignmentHandler struct {
	validate          validator.Validate
	assignmentService service.AssignmentService
}

func NewAssignmentHandler(validate validator.Validate, assignmentService service.AssignmentService) AssignmentHandler {
	return AssignmentHandler{
		validate:          validate,
		assignmentService: assignmentService,
	}
}

func (a *AssignmentHandler) CreateAssignment(c *gin.Context) {
	var assignmentRequest model.AssignmentRequest

	err := c.ShouldBindJSON(&assignmentRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = a.validate.Struct(assignmentRequest)
	if err != nil {
		logrus.Warn("Struct validation failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// TODO: Change this to the actual user from the authentication header
	assignment := assignmentRequest.ToAssignment("test")

	if err := a.assignmentService.CreateAssignment(&assignment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	assignmentResponse := assignment.ToAssignmentResponse()

	c.JSON(http.StatusCreated, assignmentResponse)

}
