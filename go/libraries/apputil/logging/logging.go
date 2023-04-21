package logging

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type CorrelationIdKey struct{}

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
	Json  bool     `yaml:"json"`
	Level LogLevel `yaml:"level"`
}

type Logger struct {
	zerolog.Logger
}

func NewLogger(config LoggerConfig) Logger {
	var writer io.Writer
	if config.Json {
		writer = os.Stdout
	} else {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	zerolog.SetGlobalLevel(config.Level.toZeroLogLevel())
	zeroLogger := zerolog.New(writer).With().Timestamp().Logger()
	return Logger{
		Logger: zeroLogger,
	}
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
	return fmt.Errorf("log level '%v' is invalid", parsedLogLevel)
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
	logger Logger
}

func NewLoggerWrapper(logger Logger) *LoggerWrapper {
	return &LoggerWrapper{logger: logger}
}

func (loggerWrapper *LoggerWrapper) Write(bytes []byte) (n int, err error) {
	loggerWrapper.logger.Error().Err(createErrorFromBytes(bytes)).Msg("server error")
	return len(bytes), nil
}

func (loggerWrapper *LoggerWrapper) ToLogger() *log.Logger {
	return log.New(loggerWrapper, "", 0)
}

func createErrorFromBytes(bytes []byte) error {
	errorString := string(bytes)
	errorString = strings.TrimSpace(errorString)
	return errors.New(errorString)
}
