package logger

import (
	"conformity-core/config"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

var Internal *logrus.Logger = New()

func New() *logrus.Logger {
	log := logrus.New()

	if config.LoggerLevel == "" {
		return log
	}

	v, err := strconv.Atoi(config.LoggerLevel)
	if err != nil {
		sentry.CaptureException(err)
		return log
	}

	log.SetLevel(logrus.Level(v))
	return log
}
