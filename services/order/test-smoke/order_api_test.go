package orderapi_test

import (
	"crypto/tls"
	"io"
	"monorepo/libraries/testutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var config = testutil.LoadConfig("../config")

func initClient() http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return http.Client{Transport: tr}
}

func TestHealth(t *testing.T) {
	// given
	client := initClient()

	// when
	response, err := client.Get(config.BaseURL + "/status/health")
	if err != nil {
		t.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(response.Body)

	// then
	assert.Equal(t, 200, response.StatusCode)
}

func TestPrometheusMetrics(t *testing.T) {
	// given
	client := initClient()

	// when
	response, err := client.Get(config.BaseURL + "/status/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(response.Body)

	// then
	assert.Equal(t, 200, response.StatusCode)
}

func TestSwaggerUI(t *testing.T) {
	// given
	client := initClient()

	// when
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

	// then
	assert.Equal(t, 200, response.StatusCode)
}
