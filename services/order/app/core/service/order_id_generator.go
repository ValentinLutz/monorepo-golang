package service

import (
	"encoding/base64"
	"fmt"
	"monorepo/libraries/apputil/config"
	"monorepo/services/order/app/core/model"
	"strings"

	"github.com/google/uuid"
)

func NewOrderId(region config.Region) model.OrderId {
	uuidValue := uuid.New()
	uuidBytes := uuidValue[:]

	base64String := base64.RawURLEncoding.EncodeToString(uuidBytes)
	base64WithoutUnderscore := replaceHyphenAndUnderscore(base64String)

	regionIdentifier := fmt.Sprintf("-%s-", region)

	base64StringHalfLength := len(base64WithoutUnderscore) / 2
	return model.OrderId(base64WithoutUnderscore[0:base64StringHalfLength] + regionIdentifier + base64WithoutUnderscore[base64StringHalfLength:])
}

func replaceHyphenAndUnderscore(input string) string {
	withoutHyphen := strings.ReplaceAll(input, "-", "!")
	withoutHyphenAndUnderscore := strings.ReplaceAll(withoutHyphen, "_", "*")
	return withoutHyphenAndUnderscore
}
