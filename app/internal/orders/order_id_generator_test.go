package orders_test

import (
	"app/internal/config"
	"app/internal/orders"
	"testing"
	"time"
)

var generatorOrderId = func(region config.Region, environment config.Environment, timestamp time.Time, salt string, expected orders.OrderId) func(t *testing.T) {
	return func(t *testing.T) {
		// WHEN
		actual := orders.GenerateOrderId(region, environment, timestamp, salt)
		// THEN
		if expected != actual {
			t.Fatalf("OrderId \n Expected: %s \n Actual: %s", expected, actual)
		}
	}
}

func Test_GenerateOrderId(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "1980-01-01T00:00:00+00:00")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("NONE-DEV", generatorOrderId(config.NONE, config.DEV, timestamp, "1", "IsQah2TkaqS-NONE-DEV-JewgL0Ye73g"))
	t.Run("NONE-DEV", generatorOrderId(config.NONE, config.DEV, timestamp, "101", "Fs2VoM7ZhrK-NONE-DEV-vzTf7kaHbRA"))
	t.Run("NONE-DEV", generatorOrderId(config.NONE, config.DEV, timestamp, "10101", "sgy1K3*SXcv-NONE-DEV-eVbldUAYXnA"))
	t.Run("NONE-DEV", generatorOrderId(config.NONE, config.DEV, timestamp, "1010101", "F2P!criGu2L-NONE-DEV-fJ7bBFx1vHg"))
	t.Run("EU-TEST", generatorOrderId(config.EU, config.TEST, timestamp, "1", "Pki*J9V8zXf-EU-TEST-wwgv9DK2f3w"))
	t.Run("EU-TEST", generatorOrderId(config.EU, config.TEST, timestamp, "101", "PHD*vw*TpU0-EU-TEST-NRMm3C6TxRg"))
	t.Run("EU-TEST", generatorOrderId(config.EU, config.TEST, timestamp, "10101", "oOhPWGD*u*j-EU-TEST-04fIDS9KysA"))
	t.Run("EU-TEST", generatorOrderId(config.EU, config.TEST, timestamp, "1010101", "gKGx34jQrzu-EU-TEST-F81CHRuomrQ"))
	t.Run("EU-PROD", generatorOrderId(config.EU, config.PROD, timestamp, "1", "Yhqp5rnLrvY-EU-9YNhvFHZLUw"))
}
