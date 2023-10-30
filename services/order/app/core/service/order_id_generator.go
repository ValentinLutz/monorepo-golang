package service

import (
	"fmt"
	"monorepo/libraries/apputil/config"
	"monorepo/services/order/app/core/model"

	"github.com/oklog/ulid/v2"
)

func NewOrderId(region config.Region) model.OrderId {
	ulidString := ulid.Make().String()
	regionIdentifier := fmt.Sprintf("-%s-", region)

	uildHalfLength := len(ulidString) / 2
	return model.OrderId(ulidString[0:uildHalfLength] + regionIdentifier + ulidString[uildHalfLength:])
}
