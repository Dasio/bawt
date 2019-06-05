package bawt

import (
	"github.com/sirupsen/logrus"
)

// Logging contains the configuration for logrus
type Logging struct {
	Logger *logrus.Logger
	Level  string `json:"level" mapstructure:"level"`
	Type   string `json:"type" mapstructure:"type"`
}

// getLoggingConfig return the corresponding formatter and level for logging.
func getLoggingConfig(bot *Bot) (logrus.Formatter, logrus.Level) {
	var f logrus.Formatter

	switch bot.Logging.Type {
	case "json":
		f = &logrus.JSONFormatter{}
	default:
		f = &logrus.TextFormatter{}
	}

	l, err := logrus.ParseLevel(bot.Logging.Level)
	if err != nil {
		l = logrus.InfoLevel
	}
	return f, l
}

// setupLogging choose the config and setup the logging.
func (bot *Bot) setupLogging() error {
	formatter, level := getLoggingConfig(bot)
	bot.Logging.Logger = logrus.New()
	log := bot.Logging.Logger

	log.SetFormatter(formatter)
	log.SetLevel(level)
	
	return nil
}
