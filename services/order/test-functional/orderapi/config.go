package orderapi

import (
	"crypto/tls"
	"monorepo/libraries/testutil"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
)

var testConfig *testutil.Config

func GetTestConfigInstance() *testutil.Config {
	if testConfig == nil {
		testConfig = testutil.LoadConfig("../../config")
	}
	return testConfig
}

var testDatabase *testutil.Database

func GetTestDatabaseInstance() *testutil.Database {
	if testDatabase == nil {
		testDatabase = testutil.NewDatabase(GetTestConfigInstance())
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

func newOrderApiClient(config *testutil.Config) *Client {
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
