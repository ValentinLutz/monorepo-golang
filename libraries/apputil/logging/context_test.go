package logging_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"monorepo/libraries/apputil/logging"
	"testing"
)

func Test_WithValue(t *testing.T) {
	// given
	ctx := context.Background()
	correlationId := "9752e638-5d49-409a-b90b-b8fd363c44eb"

	// when
	updatedCtx := logging.WithValue(ctx, logging.CorrelationIdKey, correlationId)

	// then
	value := updatedCtx.Value(logging.CorrelationIdKey)

	assert.Equal(t, correlationId, value)
}

func Benchmark_WithValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logging.WithValue(context.Background(), logging.CorrelationIdKey, "f5c25d49-dea9-4cae-bb5c-7158ba185cad")
	}
}
