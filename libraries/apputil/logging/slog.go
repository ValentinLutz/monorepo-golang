package logging

import (
	"context"
	"log/slog"
)

type ContextHandler struct {
	handler slog.Handler
}

func NewContextHandler(handler slog.Handler) *ContextHandler {
	if logHandler, ok := handler.(*ContextHandler); ok {
		handler = logHandler.handler
	}

	return &ContextHandler{
		handler: handler,
	}
}

func (contextHandler *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return contextHandler.handler.Enabled(ctx, level)
}

func (contextHandler *ContextHandler) WithGroup(name string) slog.Handler {
	return NewContextHandler(contextHandler.handler.WithGroup(name))
}

func (contextHandler *ContextHandler) Handle(ctx context.Context, record slog.Record) error {
	for key := range slogContextKeys {
		value, ok := ctx.Value(key).(SlogContextValue)
		if !ok {
			continue
		}
		record.AddAttrs(
			slog.Attr{
				Key:   key,
				Value: slog.AnyValue(value),
			},
		)
	}
	return contextHandler.handler.Handle(ctx, record)
}

func (contextHandler *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewContextHandler(contextHandler.handler.WithAttrs(attrs))
}
