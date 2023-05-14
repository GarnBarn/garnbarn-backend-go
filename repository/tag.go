package repository

import (
	"context"

	"github.com/GarnBarn/garnbarn-backend-go/model"
	"gorm.io/gorm"
)

type Tag interface {
	GetAllTag(author string) (tags []model.Tag, err error)
	Create(tag *model.Tag) error
	Update(tag *model.Tag) error
	GetByID(id int) (*model.Tag, error)
	DeleteTag(tagID int) error
}
type tag struct {
	db                *gorm.DB
	repositoryContext context.Context
}

func NewTagRepository(db *gorm.DB, repositoryContext context.Context) Tag {
	// Migrate the db
	db.AutoMigrate(&model.Tag{})

	return &tag{
		db:                db,
		repositoryContext: repositoryContext,
	}
}

func (t *tag) GetAllTag(author string) (tags []model.Tag, err error) {
	res := t.db.Model(&tags).WithContext(t.repositoryContext).Where("author = ?", author).Find(&tags)
	if res.Error != nil {
		return tags, res.Error
	}
	return tags, nil
}

func (t *tag) GetByID(id int) (*model.Tag, error) {
	tag := model.Tag{}
	result := t.db.WithContext(t.repositoryContext).First(&tag, id)
	return &tag, result.Error
}

func (t *tag) Create(tag *model.Tag) error {
	result := t.db.WithContext(t.repositoryContext).Create(tag)
	if result.Error != nil {
		return result.Error
	}

	newTag, err := t.GetByID(int(tag.ID))
	if err != nil {
		return err
	}

	*tag = *newTag
	return nil
}

func (t *tag) Update(tag *model.Tag) error {
	result := t.db.WithContext(t.repositoryContext).Save(tag)
	return result.Error
}
func (t *tag) DeleteTag(tagID int) error {
	result := t.db.WithContext(t.repositoryContext).Delete(&model.Tag{}, tagID)
	return result.Error
}
