package logging

import (
	"context"
	"golang.org/x/exp/slog"
)

type CorrelationIdKey struct{}

type ContextHandler struct {
	handler slog.Handler
	keys    map[any]string
}

func (contextHandler *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return contextHandler.handler.Enabled(ctx, level)
}

func (contextHandler *ContextHandler) WithGroup(name string) slog.Handler {
	return NewContextHandler(contextHandler.handler.WithGroup(name), contextHandler.keys)
}

func NewContextHandler(handler slog.Handler, keys map[any]string) *ContextHandler {
	if logHandler, ok := handler.(*ContextHandler); ok {
		handler = logHandler.handler
	}

	return &ContextHandler{
		handler: handler,
		keys:    keys,
	}
}

func (contextHandler *ContextHandler) Handle(ctx context.Context, record slog.Record) error {
	for key, slogKey := range contextHandler.keys {
		value, ok := ctx.Value(key).(any)
		if !ok {
			continue
		}
		record.AddAttrs(slog.Attr{
			Key:   slogKey,
			Value: slog.AnyValue(value),
		})
	}
	return contextHandler.handler.Handle(ctx, record)
}

func (contextHandler *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewContextHandler(contextHandler.handler.WithAttrs(attrs), contextHandler.keys)
}
