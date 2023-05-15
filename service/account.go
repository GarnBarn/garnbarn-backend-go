package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/GarnBarn/garnbarn-backend-go/config"
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
	"github.com/sirupsen/logrus"
)

type AccountService interface {
	GetUserByUid(uid string) (account model.AccountPublic, err error)
	CheckForCompromisedPassword(hashedPassword string) (isCompromised bool, err error)
	UpdateAccountConsentByUid(uid string, consent bool) (err error)
}

type accountService struct {
	accountRepository repository.AccountRepository
	app               *firebase.App
	appConfig         config.Config
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

func (a *accountService) UpdateAccountConsentByUid(uid string, consent bool) (err error) {
	account, err := a.accountRepository.GetAccountByUid(uid)
	if err != nil {
		logrus.Error(err)
		return err
	}

	// Set the account consent
	account.Consent = consent

	// Save the updated account
	err = a.accountRepository.UpdateAccount(account)
	return err
}

func (a *accountService) CheckForCompromisedPassword(hashedPassword string) (isCompromised bool, err error) {

	// Make request to hibp.
	prefixHashPassword := hashedPassword[0:5]
	suffixHashPassword := strings.ToUpper(hashedPassword[5:])
	url := fmt.Sprint("https://api.pwnedpasswords.com/range/", prefixHashPassword)
	logrus.Debug(url)
	resp, err := http.Get(url)
	if err != nil {
		logrus.Error(err)
		return false, err
	}

	// Validate the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return false, err
	}

	bodyString := string(body)
	bodyStringList := strings.Split(strings.ReplaceAll(bodyString, "\r\n", "\n"), "\n")

	for _, item := range bodyStringList {
		itemSplitted := strings.Split(item, ":")
		if itemSplitted[0] == suffixHashPassword {
			isCompromised = true
			break
		}
	}

	// Return the result.
	return isCompromised, nil
}
