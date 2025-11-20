package utils

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	// 配置 Logrus
	writer, _ := rotatelogs.New(
		"./logs/gin-%Y%m%d.log",
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)

	logger = logrus.New()
	logger.SetOutput(writer)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func GetLogger() *logrus.Logger {
	return logger
}

func LogRecover() {
	err := recover()
	if err != nil {
		logger.Error(err)
	}
}
