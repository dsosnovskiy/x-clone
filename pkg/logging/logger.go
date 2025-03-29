package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Init(env string) *logrus.Logger {
	logger := logrus.New()

	switch env {
	case "local":
		logger.SetLevel(logrus.DebugLevel)

		logger.SetFormatter(&logrus.TextFormatter{
			DisableLevelTruncation: true,
			ForceColors:            true,
			FullTimestamp:          true,
		})
	case "dev":
		logger.SetLevel(logrus.InfoLevel)

		logger.SetFormatter(&logrus.TextFormatter{
			DisableLevelTruncation: true,
			ForceColors:            true,
			FullTimestamp:          true,
		})
	case "prod":
		logger.SetLevel(logrus.WarnLevel)

		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	logger.SetOutput(os.Stdout)

	return logger
}
