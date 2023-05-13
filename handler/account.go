package handler

import (
	"net/http"

	"github.com/GarnBarn/garnbarn-backend-go/config"
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountHandler struct {
	accountService service.AccountService
	appConfig      config.Config
}

func NewAccountHandler(accountService service.AccountService, appConfig config.Config) AccountHandler {
	return AccountHandler{
		accountService: accountService,
		appConfig:      appConfig,
	}
}

func (a *AccountHandler) GetAccount(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		uid = c.Param(UserUidKey)
	}

	account, err := a.accountService.GetUserByUid(uid)
	if err != nil {
		logrus.Error(err)
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (a *AccountHandler) CheckForComprimizedPassword(c *gin.Context) {
	var request model.CheckCompromisedPasswordRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request Body."})
		return
	}

	if a.appConfig.HIBP_API_KEY == "" {
		// For Internal testing purpose.
		c.JSON(http.StatusOK, gin.H{"message": "Skipped check, (Internal)"})
		return
	}

	if len(request.HashedPassword) != 40 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request Body."})
		return
	}

	isCompromised, err := a.accountService.CheckForCompromisedPassword(request.HashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	if isCompromised {
		c.JSON(http.StatusFound, gin.H{"message": "Your password has been compromised"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "No passwd compromisation has been founded."})
}
