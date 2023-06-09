package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	UserUidKey = "userUid"
)

func Authentication(app *firebase.App, accountRepository repository.AccountRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("authorization")

		authHeaderSplitted := strings.Split(authHeader, " ")

		if len(authHeaderSplitted) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"messsage": "Unauthorized"})
			return
		}

		firebaseIdToken := authHeaderSplitted[1]
		ctx := context.Background()

		authClient, err := app.Auth(ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"messsage": fmt.Sprint("Middleware Error: ", err.Error())})
			return
		}

		validatedIdToken, err := authClient.VerifyIDToken(ctx, firebaseIdToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"messsage": err.Error()})
			return
		}

		account, err := accountRepository.GetAccountByUid(validatedIdToken.UID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				account, err = accountRepository.CreateAccountByUid(validatedIdToken.UID)
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"messsage": fmt.Sprint("Middleware Error: ", err.Error())})
				return
			}
		}

		if !account.Consent && c.FullPath() != "/api/v1/account/consent" {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"messsage": "Middleware Error: User doesn't accept the consent"})
			return
		}

		c.AddParam(UserUidKey, account.Uid)
		c.Next()
	}
}
