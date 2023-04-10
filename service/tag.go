package service

import (
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
)

type tag struct {
	tagRepository repository.Tag
}

type Tag interface {
	CreateTag(tag *model.Tag) error
}

func NewTagService(tagRepository repository.Tag) Tag {
	return &tag{
		tagRepository: tagRepository,
	}
}

func (t *tag) CreateTag(tag *model.Tag) error {
	return t.tagRepository.Create(tag)
}
