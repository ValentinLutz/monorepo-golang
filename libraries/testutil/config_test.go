package testutil_test

import (
	"monorepo/libraries/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseFile(t *testing.T) {
	// given
	path := "resources/test_config.yaml"

	// when
	config, err := testutil.ParseFile[testutil.Config](path)

	// then
	assert.NoError(t, err)

	expectedConfig := &testutil.Config{
		BaseURL: "http://localhost:8080",
		Database: testutil.DatabaseConfig{
			Host:     "localhost",
			Port:     9432,
			Database: "dev_db",
			Username: "test",
			Password: "test",
		},
	}
	assert.Equal(t, expectedConfig, config)
}

func Test_ParseFile_FileNotFound(t *testing.T) {
	// given
	path := "file_not_found.yaml"

	// when
	_, err := testutil.ParseFile[testutil.Config](path)

	// then
	assert.Error(t, err)
}

func Test_ParseFile_UnmarshalFailed(t *testing.T) {
	// given
	path := "resources/unknown_format.text"

	// when
	_, err := testutil.ParseFile[testutil.Config](path)

	// then
	assert.Error(t, err)
}
