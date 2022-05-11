package orders_test

import (
	"crypto/tls"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func initClient() http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return http.Client{Transport: tr}
}

func TestSwaggerUI(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	response, err := client.Get("https://localhost:8080/swagger")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	// THEN
	assert.Equal(t, 200, response.StatusCode)
}
