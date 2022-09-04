package util

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"log"
	"os"
	"strings"
	"time"
)

type LogLevel string

const (
	TRACE LogLevel = "TRACE"
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
	PANIC LogLevel = "PANIC"
)

type LoggerConfig struct {
	Pretty bool     `yaml:"pretty"`
	Level  LogLevel `yaml:"level"`
}

var logger zerolog.Logger

type Logger struct {
	logger *zerolog.Logger
}

func New() *Logger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	logger = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Logger()

	return &Logger{
		logger: &logger,
	}
}

func (logger *Logger) WithoutContext() *zerolog.Logger {
	return logger.logger
}

const CorrelationIdKey = "correlation_id"

func (logger *Logger) WithContext(context context.Context) *zerolog.Logger {
	correlationId := context.Value(CorrelationIdKey).(string)
	loggerWithContext := logger.logger.With().
		Str(CorrelationIdKey, correlationId).
		Logger()
	return &loggerWithContext
}

func SetLogLevel(logLevel LogLevel) {
	zerolog.SetGlobalLevel(logLevel.toZeroLogLevel())
}

func (logLevel *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var logLevelString string
	err := unmarshal(&logLevelString)
	if err != nil {
		return err
	}

	parsedLogLevel := LogLevel(logLevelString)

	switch parsedLogLevel {
	case TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC:
		*logLevel = parsedLogLevel
		return nil
	}
	return fmt.Errorf("log level is invalid: %s", parsedLogLevel)
}

func (logLevel *LogLevel) toZeroLogLevel() zerolog.Level {
	switch *logLevel {
	case TRACE:
		return zerolog.TraceLevel
	case DEBUG:
		return zerolog.DebugLevel
	case INFO:
		return zerolog.InfoLevel
	case WARN:
		return zerolog.WarnLevel
	case ERROR:
		return zerolog.ErrorLevel
	case FATAL:
		return zerolog.FatalLevel
	case PANIC:
		return zerolog.PanicLevel
	}
	return zerolog.NoLevel
}

type LoggerWrapper struct {
	logger *Logger
}

func NewLoggerWrapper(logger *Logger) *LoggerWrapper {
	return &LoggerWrapper{logger: logger}
}

func (lw *LoggerWrapper) Write(p []byte) (n int, err error) {
	lw.logger.logger.Error().Msg(strings.TrimSpace(string(p)))
	return len(p), nil
}

func (lw *LoggerWrapper) ToLogger() *log.Logger {
	return log.New(lw, "", 0)
}
