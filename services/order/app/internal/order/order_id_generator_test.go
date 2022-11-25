package order_test

import (
	"app/config"
	"app/core/entity"
	"app/internal/order"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testGeneratorOrderId = func(region config.Region, timestamp time.Time, salt string, expected entity.OrderId) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("Region: %v", region)
		t.Logf("Timestamp: %v", timestamp.Format(time.RFC3339))
		t.Logf("Salt: %v", salt)
		// WHEN
		actual := order.GenerateOrderId(region, timestamp, salt)
		// THEN
		assert.Equal(t, expected, actual)
	}
}

func Test_GenerateOrderId(t *testing.T) {
	// GIVEN
	timestamp, err := time.Parse(time.RFC3339, "1980-01-01T00:00:00+00:00")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("NONE", testGeneratorOrderId(config.NONE, timestamp, "1", "zanhVXdOCEg-NONE-asPc!MEMcMw"))
	t.Run("NONE", testGeneratorOrderId(config.NONE, timestamp, "101", "hjm847MUbWn-NONE-CsuoZDDc6LQ"))
	t.Run("NONE", testGeneratorOrderId(config.NONE, timestamp, "10101", "TlhDaTmRWBr-NONE-UqIiPE7q!Qw"))
	t.Run("NONE", testGeneratorOrderId(config.NONE, timestamp, "1010101", "uryHjO0*I1o-NONE-ngfDhQLBkFw"))
	t.Run("EU", testGeneratorOrderId(config.EU, timestamp, "1", "7QGZGgo5999-EU-moedOlxN4BQ"))
	t.Run("EU", testGeneratorOrderId(config.EU, timestamp, "101", "QN1iLILbclC-EU-wqVzId1oMHw"))
	t.Run("EU", testGeneratorOrderId(config.EU, timestamp, "10101", "*vFRicU14gk-EU-cA*kDJf*Jig"))
	t.Run("EU", testGeneratorOrderId(config.EU, timestamp, "1010101", "p5tCoqCnVfS-EU-J8C5J!L!mMA"))
	t.Run("US", testGeneratorOrderId(config.US, timestamp, "1", "Ad6P0F0DuUq-US-jcj2Jqrklew"))
	t.Run("US", testGeneratorOrderId(config.US, timestamp, "101", "kv0Hbli7PTn-US-TvwK!socVFg"))
	t.Run("US", testGeneratorOrderId(config.US, timestamp, "10101", "MnWZUuMf7df-US-f9GUPjtFBdA"))
	t.Run("US", testGeneratorOrderId(config.US, timestamp, "1010101", "gATm85KNU5H-US-UF!dI1xtcog"))
}
