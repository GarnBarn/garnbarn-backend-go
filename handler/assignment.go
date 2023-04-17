package handler

import (
	"net/http"
	"strconv"

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

func (a *AssignmentHandler) GetAllAssignment(c *gin.Context) {
	fromPresentString := c.Query("fromPresent")

	logrus.Debug("From Present string: ", fromPresentString)

	fromPresent := true
	if fromPresentString == "" || fromPresentString == "false" {
		fromPresent = false
	}

	logrus.Debug("From Present: ", fromPresent)

	assignments, err := a.assignmentService.GetAllAssignment(fromPresent)

	assignmentPublic := []model.AssignmentPublic{}

	for _, item := range assignments {
		assignmentPublic = append(assignmentPublic, item.ToAssignmentPublic())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.BulkResponse[model.AssignmentPublic]{
		Count:    len(assignmentPublic),
		Previous: nil,
		Next:     nil,
		Results:  assignmentPublic,
	})
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

	assignment := assignmentRequest.ToAssignment(c.Param(UserUidKey))

	if err := a.assignmentService.CreateAssignment(&assignment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	assignmentPublic := assignment.ToAssignmentPublic()

	c.JSON(http.StatusCreated, assignmentPublic)

}

func (a *AssignmentHandler) DeleteAssignment(c *gin.Context) {
	assignmentIdString, ok := c.Params.Get("Id")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}
	assignmentId, err := strconv.Atoi(assignmentIdString)
	err = a.assignmentService.DeleteAssignment(assignmentId)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}
