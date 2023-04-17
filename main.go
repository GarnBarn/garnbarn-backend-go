package main

import (
	"fmt"
	"time"

	"github.com/GarnBarn/garnbarn-backend-go/config"
	"github.com/GarnBarn/garnbarn-backend-go/handler"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/httpserver"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/logger"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-contrib/cors"
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
	assignmentRepository := repository.NewAssignmentRepository(db)

	// Create the services
	exampleService := service.NewExampleService(exampleRepository)
	tagService := service.NewTagService(tagRepository)
	assignmentService := service.NewAssignmentService(assignmentRepository)

	// Create the http server
	httpServer := httpserver.NewHttpServer()

	// Init the handler
	exampleHandler := handler.NewExampleHandler(exampleService)
	tagHandler := handler.NewTagHandler(*validate, tagService)
	assignmentHandler := handler.NewAssignmentHandler(*validate, assignmentService)

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

	// Example
	httpServer.GET("/example", exampleHandler.HelloWorld)

	// Tag
	tagRouter := router.Group("/tag")
	tagRouter.POST("/", tagHandler.CreateTag)
	tagRouter.PATCH("/:tagId", tagHandler.UpdateTag)
	tagRouter.DELETE(("/:tagId"), tagHandler.DeleteTag)

	// Assignment
	assignmentRouter := router.Group("/assignment")
	assignmentRouter.POST("/", assignmentHandler.CreateAssignment)
	assignmentRouter.DELETE("/:Id", assignmentHandler.DeleteAssignment)
	assignmentRouter.GET("/", assignmentHandler.GetAllAssignment)

	logrus.Info("Listening and serving HTTP on :", appConfig.HTTP_SERVER_PORT)
	httpServer.Run(fmt.Sprint(":", appConfig.HTTP_SERVER_PORT))
}
