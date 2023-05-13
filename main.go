package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"time"

	firebase "firebase.google.com/go"
	"github.com/GarnBarn/garnbarn-backend-go/config"
	"github.com/GarnBarn/garnbarn-backend-go/handler"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/httpserver"
	"github.com/GarnBarn/garnbarn-backend-go/pkg/logger"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-contrib/cors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	redisStore "github.com/ulule/limiter/v3/drivers/store/redis"
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

	// Initialize the Firebase App
	opt := option.WithCredentialsFile(appConfig.FIREBASE_CONFIG_FILE)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Fatalln("error initializing app: %v\n", err)
	}

	// Initialize redis
	client := redis.NewClient(&redis.Options{
		Addr:     appConfig.REDIS_CONNECTION_STRING,
		Password: appConfig.REDIS_PASSWORD,
		DB:       appConfig.REDIS_DB,
	})
	store, err := redisStore.NewStore(client)
	if err != nil {
		logrus.Fatalln("error initializing redis: %v\n", err)
	}

	// Initialize rate limiter
	rate, err := limiter.NewRateFromFormatted(appConfig.RATE_LIMIT_STRING)
	if err != nil {
		logrus.Fatalln("error initializing rate limiter: %v\n", err)
	}

	rateLimitStore := limiter.New(store, rate)
	rateLimitMiddleware := mgin.NewMiddleware(rateLimitStore)

	// Create the required dependencies
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
	assignmentHandler := handler.NewAssignmentHandler(*validate, assignmentService, tagService)
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

	router.Use(rateLimitMiddleware)
	router.Use(handler.Authentication(app, accountRepository))

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
	assignmentRouter.DELETE("/:Id", assignmentHandler.DeleteAssignment)
	assignmentRouter.GET("/", assignmentHandler.GetAllAssignment)
	assignmentRouter.GET("/:assignmentId", assignmentHandler.GetAssignmentById)
	assignmentRouter.PATCH("/:assignmentId", assignmentHandler.UpdateAssignment)

	// Account
	accountRouter := router.Group("/account")
	accountRouter.GET("/", accountHandler.GetAccount)

	logrus.Info("Listening and serving HTTP on :", appConfig.HTTP_SERVER_PORT)
	httpServer.Run(fmt.Sprint(":", appConfig.HTTP_SERVER_PORT))
}
