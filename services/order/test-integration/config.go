package testintegration

import (
	"monorepo/libraries/testutil"
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
		testDatabase = testutil.NewDatabase(GetTestConfigInstance().Database)
	}
	return testDatabase
}

var testClient *testutil.HttpClient

func GetTestClientInstance() *testutil.HttpClient {
	if testClient == nil {
		testClient = testutil.NewHttpClient(GetTestConfigInstance().BaseURL)
	}
	return testClient
}
