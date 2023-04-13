package testintegration

import (
	"monorepo/libraries/testingutil"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	config   *testingutil.Config
	database *sqlx.DB
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.config = testingutil.LoadConfig("../config")
	suite.database = testingutil.NewDatabase(suite.T(), suite.config)

	testingutil.LoadAndExec(suite.T(), suite.database, "orderapi/truncate_tables.sql")
}

func Test_IntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
