package order_api_test

import (
	"crypto/tls"
	"io"
	"monorepo/libraries/testingutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var config = testingutil.LoadConfig("../config")

func initClient() http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return http.Client{Transport: tr}
}

func TestHealth(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	response, err := client.Get(config.BaseURL + "/api/status/health")
	if err != nil {
		t.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(response.Body)

	// THEN
	assert.Equal(t, 200, response.StatusCode)
}

func TestPrometheusMetrics(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	response, err := client.Get(config.BaseURL + "/api/status/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(response.Body)

	// THEN
	assert.Equal(t, 200, response.StatusCode)
}

func TestSwaggerUI(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	response, err := client.Get(config.BaseURL + "/swagger/")
	if err != nil {
		t.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(response.Body)

	// THEN
	assert.Equal(t, 200, response.StatusCode)
}
