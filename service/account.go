package service

import (
	"context"

	firebase "firebase.google.com/go"
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
	"github.com/sirupsen/logrus"
)

type AccountService interface {
	GetUserByUid(uid string) (account model.AccountPublic, err error)
	CheckForCompromisedPassword(hashedPassword string) (isCompromised bool, err error)
}

type accountService struct {
	accountRepository repository.AccountRepository
	app               *firebase.App
}

func NewAccountService(app *firebase.App, accountRepository repository.AccountRepository) AccountService {
	return &accountService{
		app:               app,
		accountRepository: accountRepository,
	}
}

func (a *accountService) GetUserByUid(uid string) (account model.AccountPublic, err error) {
	// Get Account From Database
	accountPrivate, err := a.accountRepository.GetAccountByUid(uid)
	if err != nil {
		logrus.Error("Can't get account from database: ", err)
		return account, err
	}

	// Fill the Account Information by using data from Firebase

	ctx := context.Background()
	auth, err := a.app.Auth(ctx)
	if err != nil {
		return account, err
	}

	user, err := auth.GetUser(ctx, uid)
	if err != nil {
		return account, err
	}

	return accountPrivate.ToAccountPublic(user.DisplayName, user.PhotoURL), nil
}

func (a *accountService) CheckForCompromisedPassword(hashedPassword string) (isCompromised bool, err error) {

	// TODO: Make request to ihbp.

	// TODO: Validate the response.

	// TODO: Return the result.

	return false, nil
}
