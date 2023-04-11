package handler

import (
	"net/http"

	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Tag struct {
	validate   validator.Validate
	tagService service.Tag
}

func NewTagHandler(validate validator.Validate, tagService service.Tag) Tag {
	return Tag{
		validate:   validate,
		tagService: tagService,
	}
}

func (t *Tag) CreateTag(c *gin.Context) {

	var tagRequest model.CreateTagRequest

	err := c.ShouldBind(&tagRequest)
	if err != nil {
		logrus.Warn("Requets Body binding error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request body."})
		return
	}

	err = t.validate.Struct(tagRequest)
	if err != nil {
		logrus.Warn("Struct validation failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// TODO: Change this to the actual user from the authentication header
	tag := tagRequest.ToTag("test")

	err = t.tagService.CreateTag(&tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in the server"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
