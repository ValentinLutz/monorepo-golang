package orderapi

import (
	"crypto/tls"
	"monorepo/libraries/testingutil"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
)

var testConfig *testingutil.Config

func GetTestConfigInstance() *testingutil.Config {
	if testConfig == nil {
		testConfig = testingutil.LoadConfig("../../config")
	}
	return testConfig
}

var testDatabase *testingutil.Database

func GetTestDatabaseInstance() *testingutil.Database {
	if testDatabase == nil {
		testDatabase = testingutil.NewDatabase(GetTestConfigInstance())
	}
	return testDatabase
}

var orderApiClient *Client

func GetOrderApiClientInstance() *Client {
	if orderApiClient == nil {
		orderApiClient = newOrderApiClient(GetTestConfigInstance())
	}
	return orderApiClient
}

func newOrderApiClient(config *testingutil.Config) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	basicAuth, err := securityprovider.NewSecurityProviderBasicAuth("test", "test")
	if err != nil {
		panic(err)
	}

	orderApiClient, err := NewClient(
		config.BaseURL+"/",
		WithHTTPClient(client),
		WithRequestEditorFn(basicAuth.Intercept),
	)
	if err != nil {
		panic(err)
	}

	return orderApiClient
}
