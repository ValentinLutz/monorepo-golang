package order_test

import (
	"app/core/entity"
	"app/internal/config"
	"app/internal/order"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var generatorOrderId = func(region config.Region, environment config.Environment, timestamp time.Time, salt string, expected entity.OrderId) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("Region: %v", region)
		t.Logf("Environment: %v", environment)
		t.Logf("Timestamp: %v", timestamp.Format(time.RFC3339))
		t.Logf("Salt: %v", salt)
		// WHEN
		actual := order.GenerateOrderId(region, environment, timestamp, salt)
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
	t.Run("NONE-DEV", generatorOrderId(config.NONE, config.DEV, timestamp, "1", "fdCDxjV9o!O-NONE-DEV-ZCTH5i6fWcA"))
	t.Run("NONE-DEV", generatorOrderId(config.NONE, config.DEV, timestamp, "101", "X!d2QGwSqbz-NONE-DEV-O83j82*zsIw"))
	t.Run("NONE-DEV", generatorOrderId(config.NONE, config.DEV, timestamp, "10101", "10GQefiu13u-NONE-DEV-6DiC6mjULDw"))
	t.Run("NONE-DEV", generatorOrderId(config.NONE, config.DEV, timestamp, "1010101", "eHYnyDHx61P-NONE-DEV-xyozdTi9jcA"))
	t.Run("EU-TEST", generatorOrderId(config.EU, config.TEST, timestamp, "1", "WYwajVCTfxv-EU-TEST-R16B1EOIYqA"))
	t.Run("EU-TEST", generatorOrderId(config.EU, config.TEST, timestamp, "101", "*SOY*UcOhPS-EU-TEST-Ph7SlmnWyPA"))
	t.Run("EU-TEST", generatorOrderId(config.EU, config.TEST, timestamp, "10101", "xn*jZKcJI0e-EU-TEST-XjStk3UhSxw"))
	t.Run("EU-TEST", generatorOrderId(config.EU, config.TEST, timestamp, "1010101", "wUq!zyfXs!a-EU-TEST-Tm4spTd4IRA"))
	t.Run("EU-PROD", generatorOrderId(config.EU, config.PROD, timestamp, "1", "xnl7K9M2NUw-EU-oLpPbSQHEWA"))
}
