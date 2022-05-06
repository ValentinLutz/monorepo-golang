package orders_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"test-integration/orders"
	"testing"
)

func initClient() orders.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return orders.Client{
		Server: "https://localhost:8080",
		Client: client,
	}
}

func TestGetOrder(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	apiOrder, _ := client.GetApiOrdersOrderId(context.Background(), "IsQah2TkaqS-NONE-DEV-JewgL0Ye73g")
	defer apiOrder.Body.Close()

	// THEN
	var actualResponse orders.OrderResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orders.OrderResponse
	readToObject(t, readFile(t, "orderResponse.json"), &expectedResponse)

	assert.Equal(t, 200, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestGetOrderNotFound(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	apiOrder, _ := client.GetApiOrdersOrderId(context.Background(), "NOPE")
	defer apiOrder.Body.Close()

	// THEN
	var actualResponse orders.ErrorResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orders.ErrorResponse
	readToObject(t, readFile(t, "orderNotFoundResponse.json"), &expectedResponse)
	expectedResponse.Timestamp = actualResponse.Timestamp

	assert.Equal(t, 404, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestGetOrders(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	apiOrder, _ := client.GetApiOrders(context.Background())
	defer apiOrder.Body.Close()

	// THEN
	var actualResponse orders.OrdersResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orders.OrdersResponse
	readToObject(t, readFile(t, "ordersResponse.json"), &expectedResponse)

	assert.Equal(t, 200, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
}

func readToObject(t *testing.T, reader io.Reader, object interface{}) {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(object)
	if err != nil {
		t.Fatalf("Failed to decode input, %v", err)
	}
}

func readFile(t *testing.T, path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to read file from path %v, %v", path, err)
	}
	return file
}
