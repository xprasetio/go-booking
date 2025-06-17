package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger adalah interface untuk logging
type Logger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	Error(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Fatal(args ...interface{})
}

// NewLogger membuat instance logger baru
func NewLogger() Logger {
	log := logrus.New()

	// Set format output ke JSON
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set output ke stdout
	log.SetOutput(os.Stdout)

	// Set level log berdasarkan environment
	if os.Getenv("APP_ENV") == "production" {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logrus.DebugLevel)
	}

	return log
}
