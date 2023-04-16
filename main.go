package main

import (
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go"
	"github.com/GarnBarn/garnbarn-backend-go/config"
	"github.com/GarnBarn/garnbarn-backend-go/handler"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/httpserver"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/logger"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/middleware"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-contrib/cors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	appConfig config.Config
)

func init() {
	appConfig = config.Load()
	logger.InitLogger(logger.Config{
		Env: appConfig.Env,
	})

}

func main() {
	// Start DB Connection
	db, err := gorm.Open(mysql.Open(appConfig.MYSQL_CONNECTION_STRING), &gorm.Config{})
	if err != nil {
		logrus.Panic("Can't connect to db: ", err)
	}

	// Initilize the Firebase App
	opt := option.WithCredentialsFile(appConfig.FIREBASE_CONFIG_FILE)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Fatalln("error initializing app: %v\n", err)
	}

	// Create the required dependentices
	validate := validator.New()

	// Create the repositroies
	tagRepository := repository.NewTagRepository(db)
	assignmentRepository := repository.NewAssignmentRepository(db)
	accountRepository := repository.NewAccountRepository(db)

	// Create the services
	tagService := service.NewTagService(tagRepository)
	assignmentService := service.NewAssignmentService(assignmentRepository)
	accountService := service.NewAccountService(app, accountRepository)

	// Create the http server
	httpServer := httpserver.NewHttpServer()

	// Init the handler
	tagHandler := handler.NewTagHandler(*validate, tagService)
	assignmentHandler := handler.NewAssignmentHandler(*validate, assignmentService)
	accountHandler := handler.NewAccountHandler(accountService)

	httpServer.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// Router
	router := httpServer.Group("/api/v1")

	router.Use(middleware.Authentication(app))

	// Tag
	tagRouter := router.Group("/tag")
	tagRouter.GET("/:id", tagHandler.GetTagById)
	tagRouter.GET("/", tagHandler.GetAllTag)
	tagRouter.POST("/", tagHandler.CreateTag)
	tagRouter.PATCH("/:tagId", tagHandler.UpdateTag)
	tagRouter.DELETE("/:tagId", tagHandler.DeleteTag)

	// Assignment
	assignmentRouter := router.Group("/assignment")
	assignmentRouter.POST("/", assignmentHandler.CreateAssignment)
	assignmentRouter.GET("/", assignmentHandler.GetAllAssignment)

	// Account
	accountRouter := router.Group("/account")
	accountRouter.GET("/", accountHandler.GetAccount)

	logrus.Info("Listening and serving HTTP on :", appConfig.HTTP_SERVER_PORT)
	httpServer.Run(fmt.Sprint(":", appConfig.HTTP_SERVER_PORT))
}
