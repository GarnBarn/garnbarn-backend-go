package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/GarnBarn/garnbarn-backend-go/config"
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountHandler struct {
	accountService service.AccountService
	appConfig      config.Config
	app            *firebase.App
}

func NewAccountHandler(accountService service.AccountService, appConfig config.Config, app *firebase.App) AccountHandler {
	return AccountHandler{
		accountService: accountService,
		appConfig:      appConfig,
		app:            app,
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
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in server."})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (a *AccountHandler) UpdateAccountConsentByUid(c *gin.Context) {
	// Bind the request body.
	var updateAccountRequest model.UpdateAccountRequest
	err := c.ShouldBindJSON(&updateAccountRequest)
	if err != nil {
		logrus.Warn("Can't bind request body to model: ", err)
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	err = a.accountService.UpdateAccountConsentByUid(c.Param(UserUidKey), updateAccountRequest.Consent)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in the server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
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
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in server."})
		return
	}

	if isCompromised {
		c.JSON(http.StatusFound, gin.H{"message": "Your password has been compromised"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "No passwd compromisation has been founded."})
}

func (a *AccountHandler) GetToken(c *gin.Context) {
	var request model.GetTokenRequest
	err := c.ShouldBind(&request)
	if err != nil || request.Code == "" || request.RedirectUri == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request Body."})
		return
	}

	ctx := context.Background()
	authClient, err := a.app.Auth(ctx)

	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happened in server."})
		return
	}

	// Send to get token
	data := url.Values{}
	data.Set("code", request.Code)
	data.Set("redirect_uri", request.RedirectUri)
	data.Set("client_id", a.appConfig.OAUTH_CLIENT_ID)
	data.Set("client_secret", a.appConfig.OAUTH_CLIENT_SECRET)
	data.Set("grant_type", "authorization_code")
	encodedData := data.Encode()

	resp, err := http.Post("https://garnbarn.jp.auth0.com/oauth/token", "application/x-www-form-urlencoded", strings.NewReader(encodedData))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var parsedTokenBody model.OAuthToken
	err = json.Unmarshal(respBody, &parsedTokenBody)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	logrus.Debug(parsedTokenBody.IDToken)

	token, _, err := new(jwt.Parser).ParseUnverified(parsedTokenBody.IDToken, jwt.MapClaims{})
	if err != nil {
		logrus.Error("JWT Error: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	email := claims["email"].(string)
	nonce := claims["nonce"].(string)

	if nonce != request.Nonce {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	user, err := authClient.GetUserByEmail(ctx, email)
	if err != nil {
		logrus.Warn("error getting user by email %s: %v\n", email, err)
		newUser := auth.UserToCreate{}
		newUser.Email(email)
		user, err = authClient.CreateUser(ctx, &newUser)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
	}

	customToken, err := authClient.CustomToken(ctx, user.UID)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in server"})
		return
	}

	c.JSON(http.StatusOK, model.TokenResponse{
		Token: customToken,
	})

}
