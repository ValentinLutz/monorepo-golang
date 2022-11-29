package order

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"monorepo/services/order/app/config"
	"monorepo/services/order/app/core/entity"
	"strings"
	"time"
)

func NewOrderId(region config.Region, timestamp time.Time, salt string) entity.OrderId {
	valueToHash := string(region) + timestamp.Format(time.RFC3339) + salt
	md5Sum := md5.Sum([]byte(valueToHash))

	base64String := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(md5Sum[:])
	base64WithoutUnderscore := replaceHyphenAndUnderscore(base64String)

	regionIdentifier := fmt.Sprintf("-%s-", region)

	base64StringHalfLength := len(base64WithoutUnderscore) / 2
	return entity.OrderId(base64WithoutUnderscore[0:base64StringHalfLength] + regionIdentifier + base64WithoutUnderscore[base64StringHalfLength:])
}

func replaceHyphenAndUnderscore(input string) string {
	withoutHyphen := strings.ReplaceAll(input, "-", "!")
	withoutHyphenAndUnderscore := strings.ReplaceAll(withoutHyphen, "_", "*")
	return withoutHyphenAndUnderscore
}
