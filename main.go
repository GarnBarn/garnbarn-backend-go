package main

import (
	"fmt"

	"github.com/GarnBarn/garnbarn-backend-go/config"
	"github.com/GarnBarn/garnbarn-backend-go/handler"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/httpserver"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/logger"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
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

	// Create the required dependentices
	validate := validator.New()

	// Create the repositroies
	exampleRepository := repository.NewExampleRepository(db)
	tagRepository := repository.NewTagRepository(db)

	// Create the services
	exampleService := service.NewExampleService(exampleRepository)
	tagService := service.NewTagService(tagRepository)

	// Create the http server
	httpServer := httpserver.NewHttpServer()

	// Init the handler
	exampleHandler := handler.NewExampleHandler(exampleService)
	tagHandler := handler.NewTagHandler(*validate, tagService)

	// Add Routes
	httpServer.GET("/example", exampleHandler.HelloWorld)

	// Tag
	tagRouter := httpServer.Group("/tag")
	tagRouter.GET("/", tagHandler.GetAllTags)
	tagRouter.GET("/:id", tagHandler.GetTagById)
	tagRouter.POST("/", tagHandler.CreateTag)
	tagRouter.PATCH("/:tagId", tagHandler.UpdateTag)

	logrus.Info("Listening and serving HTTP on :", appConfig.HTTP_SERVER_PORT)
	httpServer.Run(fmt.Sprint(":", appConfig.HTTP_SERVER_PORT))
}
