package order_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"test-integration/order"
	"testing"
)

func initClient() order.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return order.Client{
		Server: "http://localhost:8080",
		Client: client,
	}
}

func TestPostOrder(t *testing.T) {
	// GIVEN
	client := initClient()
	orderItems := []order.OrderItemRequest{
		{Name: "caramel"},
		{Name: "clementine"},
	}
	orderRequest := order.OrderRequest{Items: orderItems}
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(orderRequest)
	if err != nil {
		t.Fatal(err)
	}

	// WHEN
	apiOrder, err := client.PostApiOrdersWithBody(context.Background(), "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	var actualResponse order.OrderResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	assert.Equal(t, 201, apiOrder.StatusCode)
	assert.Equal(t, order.OrderPlaced, actualResponse.Status)
	assert.Equal(t, []order.OrderItemResponse{
		{Name: "caramel"},
		{Name: "clementine"},
	}, actualResponse.Items)
	assert.NotEmpty(t, actualResponse.OrderId)
	assert.NotEmpty(t, actualResponse.CreationDate)
}

func TestGetOrder(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	apiOrder, err := client.GetApiOrdersOrderId(context.Background(), "IsQah2TkaqS-NONE-DEV-JewgL0Ye73g")
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	var actualResponse order.OrderResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse order.OrderResponse
	readToObject(t, readFile(t, "orderResponse.json"), &expectedResponse)

	assert.Equal(t, 200, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestGetOrderNotFound(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	apiOrder, err := client.GetApiOrdersOrderId(context.Background(), "NOPE")
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	var actualResponse order.ErrorResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse order.ErrorResponse
	readToObject(t, readFile(t, "orderNotFoundResponse.json"), &expectedResponse)
	expectedResponse.Timestamp = actualResponse.Timestamp

	assert.Equal(t, 404, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestGetOrders(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	apiOrder, err := client.GetApiOrders(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	var actualResponse order.OrdersResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse order.OrdersResponse
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
