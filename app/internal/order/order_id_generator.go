package order

import (
	"app/internal/config"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

func GenerateOrderId(region config.Region, environment config.Environment, timestamp time.Time, salt string) OrderId {
	valueToHash := string(region) + string(environment) + timestamp.Format(time.RFC3339) + salt
	md5Sum := md5.Sum([]byte(valueToHash))

	base64String := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(md5Sum[:])
	base64WithoutUnderscore := replaceHyphenAndUnderscore(base64String)

	regionIdentifier := buildRegionIdentifier(region, environment)

	base64StringHalfLength := len(base64WithoutUnderscore) / 2
	return OrderId(base64WithoutUnderscore[0:base64StringHalfLength] + regionIdentifier + base64WithoutUnderscore[base64StringHalfLength:])
}

func replaceHyphenAndUnderscore(input string) string {
	withoutHyphen := strings.ReplaceAll(input, "-", "!")
	withoutHyphenAndUnderscore := strings.ReplaceAll(withoutHyphen, "_", "*")
	return withoutHyphenAndUnderscore
}

func buildRegionIdentifier(region config.Region, environment config.Environment) string {
	if environment == config.PROD {
		return fmt.Sprintf("-%s-", region)
	}
	return fmt.Sprintf("-%s-%s-", region, environment)
}
