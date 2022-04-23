package internal

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

type LoggerConfig struct {
	Pretty bool   `yaml:"pretty"`
	Level  string `yaml:"level"`
}

func NewLogger(level zerolog.Level, prettyLogging bool) zerolog.Logger {
	var out io.Writer
	out = os.Stderr

	if prettyLogging {
		out = zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: time.RFC3339,
		}
	}

	logger := zerolog.New(out).
		With().
		Timestamp().
		Logger()
	logger.WithLevel(level)

	return logger
}
