package handler

import (
	"net/http"

	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-gonic/gin"
)

type Tag struct {
	tagService service.Tag
}

func NewTagHandler(tagService service.Tag) Tag {
	return Tag{
		tagService: tagService,
	}
}

func (t *Tag) CreateTag(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
