package service

import (
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
	"github.com/pquerna/otp/totp"
	"github.com/sirupsen/logrus"
)

type tag struct {
	tagRepository repository.Tag
}

type Tag interface {
	CreateTag(tag *model.Tag) error
	UpdateTag(tagId int, tagUpdateRequest *model.UpdateTagRequest) (*model.Tag, error)
	GetTagById(tagId string) (model.TagPublic, error)
}

func NewTagService(tagRepository repository.Tag) Tag {
	return &tag{
		tagRepository: tagRepository,
	}
}

func (t *tag) CreateTag(tag *model.Tag) error {

	// Create the otp secret
	totpKeyResult, err := totp.Generate(totp.GenerateOpts{Issuer: "GarnBarn", AccountName: "GarnBarn"})
	if err != nil {
		logrus.Error(err)
		return err
	}
	totpPrivateKey := totpKeyResult.Secret()
	logrus.Info(totpPrivateKey)

	tag.SecretKeyTotp = totpPrivateKey

	return t.tagRepository.Create(tag)
}

func (t *tag) UpdateTag(tagId int, tagUpdateRequest *model.UpdateTagRequest) (*model.Tag, error) {
	// Get current tag
	tag, err := t.tagRepository.GetByID(tagId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Update the tag object
	tagUpdateRequest.UpdateTag(tag)

	// Update the data in db.
	err = t.tagRepository.Update(tag)
	return tag, err
}

func (t *tag) GetTagById(tagId string) (model.TagPublic, error) {
	tag, err := t.tagRepository.GetTagById(tagId)
	if err != nil {
		logrus.Error(err)
		return model.TagPublic{}, err
	}

	return tag.ToTagPublic(), nil
}
