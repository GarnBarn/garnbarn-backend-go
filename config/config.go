package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Env                      string
	HTTP_SERVER_PORT         string   `envconfig:"HTTP_SERVER_PORT" default:"3000"`
	GIN_MODE                 string   `envconfig:"GIN_MODE" default:"release"`
	MYSQL_CONNECTION_STRING  string   `envconfig:"MYSQL_CONNECTION_STRING"`
	FIREBASE_CONFIG_FILE     string   `envconfig:"FIREBASE_CONFIG_FILE" default:"firebase-credential.json"`
	ALLOW_ORIGINS            []string `envconfig:"ALLOW_ORIGINS" default:"*"`
	ALLOW_METHODS            []string `envconfig:"ALLOW_METHODS" default:"*"`
	ALLOW_HEADERS            []string `envconfig:"ALLOW_HEADERS" default:"*"`
	EXPOSE_HEADERS           []string `envconfig:"EXPOSE_HEADERS" default:"*"`
	ALLOW_CREDENTIALS        bool     `envconfig:"ALLOW_CREDENTIALS" default:"true"`
	MAX_AGE                  int      `envconfig:"MAX_AGE" default:"12"`
	REDIS_CONNECTION_STRING  string   `envconfig:"REDIS_CONNECTION_STRING"`
	REDIS_PASSWORD           string   `envconfig:"REDIS_PASSWORD"`
	RATE_LIMIT_STRING        string   `envconfig:"RATE_LIMIT_STRING" default:"5-S"`
	REDIS_DB                 int      `envconfig:"REDIS_DB" default:"0"`
	HIBP_API_KEY             string   `envconfig:"HIBP_API_KEY"`
	TAG_TABLE_AES_KEY        string   `envconfig:"TAG_TABLE_AES_KEY"`
	ASSIGNMENT_TABLE_AES_KEY string   `envconfig:"ASSIGNMENT_TABLE_AES_KEY"`
	ACCOUNT_TABLE_AES_KEY    string   `envconfig:"ACCOUNT_TABLE_AES_KEY"`
}

const (
	TagEncryptionContextKey        string = "tagEncryptionKey"
	AssignmentEncryptionContextKey string = "assignmentEncryptionKey"
)

func Load() Config {
	var config Config
	ENV, ok := os.LookupEnv("ENV")
	if !ok {
		// Default value for ENV.
		ENV = "dev"
	}

	if ENV == "prod" {
		timeFormatLayout := "2006-01-02T15:04:05.000Z"
		logrus.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "log_level",
			},
			TimestampFormat: timeFormatLayout,
		})
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.TraceLevel)
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	// Load the .env file only for dev env.
	ENV_CONFIG, ok := os.LookupEnv("ENV_CONFIG")
	if !ok {
		ENV_CONFIG = "./.env"
	}

	err := godotenv.Load(ENV_CONFIG)
	if err != nil {
		logrus.Warn("Can't load env file")
	}

	envconfig.MustProcess("", &config)
	config.Env = ENV

	return config
}
