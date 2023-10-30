package service_test

import (
	"monorepo/libraries/apputil/config"
	"monorepo/services/order/app/core/service"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Benchmark_NewOrderId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		service.NewOrderId(config.NONE)
	}
}

func Test_NewOrderId(t *testing.T) {
	// given
	regions := []string{
		"NONE",
		"EU",
		"US",
	}
	regex := regexp.MustCompile("^[A-Za-z0-9]{13}-[A-Z]{2,4}-[A-Za-z0-9]{13}$")

	for _, region := range regions {
		for i := 0; i < 100; i++ {
			t.Run(region, testNewOrderId(config.Region(region), regex))
		}
	}
}

func testNewOrderId(region config.Region, regex *regexp.Regexp) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("Region: %v", region)

		// when
		orderId := service.NewOrderId(region)
		// then
		assert.Regexp(t, regex, orderId)
	}
}
