package orders_test

import (
	"app/internal/config"
	"app/internal/orders"
	"testing"
)

var generatorOrderId = func(region config.Region, environment config.Environment, salt string, expected orders.OrderId) func(t *testing.T) {
	return func(t *testing.T) {
		// WHEN
		actual := orders.GenerateOrderId(region, environment, salt)
		// THEN
		if expected != actual {
			t.Fatalf("OrderId \n Expected: %s \n Actual: %s", expected, actual)
		}
	}
}

func Test_GenerateOrderId(t *testing.T) {
	t.Run("NONE-DEV", generatorOrderId(config.NONE, config.DEV, "1", "80g40Ma!Jp5-NONE-DEV-R6JetvWObxw"))
	t.Run("EU-TEST", generatorOrderId(config.EU, config.TEST, "101", "6PMVaGPemWu-EU-TEST-WfNutgjYl8A"))
	t.Run("EU-PROD", generatorOrderId(config.EU, config.PROD, "10101", "mmsmDgqpMzS-EU-VQA6IK3oUBQ"))
}
