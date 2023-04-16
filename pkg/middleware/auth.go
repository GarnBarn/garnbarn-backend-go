package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
)

func Authentication(app *firebase.App) gin.HandlerFunc {
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

		_, err = authClient.VerifyIDToken(ctx, firebaseIdToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"messsage": err.Error()})
			return
		}

	}
}
