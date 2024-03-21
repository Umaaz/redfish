package logging

import (
	"log/slog"
	"os"
	"strings"
)

var (
	Logger   *slog.Logger
	logLevel *slog.LevelVar
)

func init() {
	println("init logging")
	// setup logging
	logLevel = new(slog.LevelVar)
	logLevel.Set(LevelFromEnv())
	if os.Getenv("LOG_TYPE") == "json" {
		Logger = JsonLogging(logLevel)
	} else {
		Logger = TextLogging(logLevel)
	}
}

// JsonLogging will create a new slog handler that will print logs to stderr as json
func JsonLogging(lvl *slog.LevelVar) *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	})
	return slog.New(jsonHandler)
}

// TextLogging will create a new slog handler that will print logs to stderr as json
func TextLogging(lvl *slog.LevelVar) *slog.Logger {
	jsonHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	})
	return slog.New(jsonHandler)
}

// LevelFromEnv will read the LOG_LEVEL environment variable to set the initial log level
func LevelFromEnv() slog.Level {
	logLevel := os.Getenv("LOG_LEVEL")
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
	case "":
	case "INFO":
		return slog.LevelInfo
	}
	return slog.LevelInfo
}
