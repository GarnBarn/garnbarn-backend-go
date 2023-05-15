package repository

import (
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"gorm.io/gorm"
)

type AccountRepository interface {
	GetAccountByUid(uid string) (account model.Account, err error)
	CreateAccountByUid(uid string) (account model.Account, err error)
	UpdateAccount(account model.Account) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	db.AutoMigrate(&model.Account{})

	return &accountRepository{
		db: db,
	}
}

func (a *accountRepository) GetAccountByUid(uid string) (account model.Account, err error) {
	res := a.db.First(&account, "uid = ?", uid)
	if res.Error != nil {
		return account, res.Error
	}
	return account, nil
}

func (a *accountRepository) CreateAccountByUid(uid string) (account model.Account, err error) {
	account = model.Account{
		Uid:     uid,
		Consent: false,
	}
	res := a.db.Save(&account)
	if res.Error != nil {
		return account, res.Error
	}
	return account, err
}

func (a *accountRepository) UpdateAccount(account model.Account) error {
	result := a.db.Save(&account)
	return result.Error
}
