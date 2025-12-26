package utils

import (
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var multiWriter io.Writer

func init() {
	// 创建日志目录
	if err := os.MkdirAll("./logs", 0755); err != nil {
		panic(err)
	}

	// 配置轮转日志 writer
	fileWriter, err := rotatelogs.New(
		"./logs/gin-%Y%m%d.log",
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)
	if err != nil {
		panic(err)
	}

	// 创建多写入器，同时写入文件和控制台
	multiWriter = io.MultiWriter(fileWriter, os.Stdout)

	logger = logrus.New()
	logger.SetOutput(multiWriter)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

func GetMultiWriter() io.Writer {
	return multiWriter
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
