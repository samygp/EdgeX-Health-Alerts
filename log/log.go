package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/samygp/edgex-health-alerts/config"
	"github.com/samygp/edgex-health-alerts/version"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents the Logger instance to be used.
var Logger *zap.SugaredLogger

//Init the logger
func Init() {
	var conf zap.Config

	if config.Config.App.Debug {
		conf = zap.NewDevelopmentConfig()
		conf.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	} else {
		conf = zap.NewProductionConfig()
		conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	conf.Level = zap.NewAtomicLevelAt(parseLoggerLevel(config.Config.Logger.Level))
	logger, err := conf.Build()
	if err != nil {
		fmt.Printf("Unable to create logger: %v", err)
	} else {
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Printf("Unable to get hostname: %v", err)
		}

		Logger = logger.Named(version.Name).With(zap.String("name", version.ID), zap.String("hostname", hostname)).Sugar()
	}
}

func parseLoggerLevel(level string) zapcore.Level {
	lvl := zap.FatalLevel

	if strings.EqualFold(level, "DEBUG") {
		lvl = zap.DebugLevel
	} else if strings.EqualFold(level, "INFO") {
		lvl = zap.InfoLevel
	} else if strings.EqualFold(level, "WARN") {
		lvl = zap.WarnLevel
	} else if strings.EqualFold(level, "ERROR") {
		lvl = zap.ErrorLevel
	}

	return lvl
}
