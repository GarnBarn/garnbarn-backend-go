package logger

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Env string
}

func InitLogger(loggerConfig Config) {
	if loggerConfig.Env == "prod" {
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
}
