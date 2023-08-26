package logging

import (
	"context"
)

var CorrelationIdKey = SlogContextKey{Name: "correlation_id"}

type SlogContextKey struct {
	Name string
}

type SlogContextValue any

var slogContextKeys = make(map[string]any)

func WithValue(ctx context.Context, key SlogContextKey, value SlogContextValue) context.Context {
	_, ok := slogContextKeys[key.Name]
	if !ok {
		slogContextKeys[key.Name] = nil
	}
	return context.WithValue(ctx, key, value)
}
