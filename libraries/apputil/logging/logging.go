package logging

import (
	"fmt"
	"log"
	"log/slog"
	"os"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

type LoggerConfig struct {
	Level LogLevel `yaml:"level"`
}

func (logLevel *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var logLevelString string
	err := unmarshal(&logLevelString)
	if err != nil {
		return err
	}

	parsedLogLevel := LogLevel(logLevelString)

	switch parsedLogLevel {
	case DEBUG, INFO, WARN, ERROR:
		*logLevel = parsedLogLevel
		return nil
	}
	return fmt.Errorf("invalid log level '%v'", parsedLogLevel)
}

func (logLevel *LogLevel) toSlogLevel() (slog.Level, error) {
	switch *logLevel {
	case DEBUG:
		return slog.LevelDebug, nil
	case INFO:
		return slog.LevelInfo, nil
	case WARN:
		return slog.LevelWarn, nil
	case ERROR:
		return slog.LevelError, nil
	}
	return 0, fmt.Errorf("unknown log level '%v'", *logLevel)
}

func NewSlogHandler(config LoggerConfig) slog.Handler {
	slogLevel, err := config.Level.toSlogLevel()
	if err != nil {
		slog.Error(
			"failed to parse log level",
			slog.Any("err", err),
		)
		os.Exit(1)
	}

	return slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slogLevel,
		},
	)
}

func NewSlogLogger(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}

func NewLogger(handler slog.Handler, config LoggerConfig) *log.Logger {
	slogLevel, err := config.Level.toSlogLevel()
	if err != nil {
		slog.Error(
			"failed to parse log level",
			slog.Any("err", err),
		)
		os.Exit(1)
	}

	return slog.NewLogLogger(handler, slogLevel)
}
