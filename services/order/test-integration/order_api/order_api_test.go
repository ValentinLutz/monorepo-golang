package order_api_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/ValentinLutz/monrepo/libraries/testingutil"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"test-integration/order_api"
	"testing"
	"time"
)

var config = testingutil.LoadConfig("../../config/test")

func initClient() order_api.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return order_api.Client{
		Server: config.BaseURL,
		Client: client,
	}
}

func TestGetOrders(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	startTime := time.Now()
	apiOrder, err := client.GetApiOrders(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()
	responseTimeInMs := time.Since(startTime).Milliseconds()

	// THEN
	var actualResponse order_api.OrdersResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse order_api.OrdersResponse
	readToObject(t, readFile(t, "ordersResponse.json"), &expectedResponse)

	assert.Equal(t, 200, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
	assert.GreaterOrEqual(t, int64(10), responseTimeInMs, "Response time in milliseconds")
}

func TestPostOrder(t *testing.T) {
	// GIVEN
	client := initClient()
	orderItems := []order_api.OrderItemRequest{
		{Name: "caramel"},
		{Name: "clementine"},
	}
	orderRequest := order_api.OrderRequest{Items: orderItems}
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(orderRequest)
	if err != nil {
		t.Fatal(err)
	}

	// WHEN
	startTime := time.Now()
	apiOrder, err := client.PostApiOrdersWithBody(context.Background(), "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()
	responseTimeInMs := time.Since(startTime).Milliseconds()

	// THEN
	var actualResponse order_api.OrderResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	assert.Equal(t, 201, apiOrder.StatusCode)
	assert.Equal(t, order_api.OrderPlaced, actualResponse.Status)
	assert.Equal(t, []order_api.OrderItemResponse{
		{Name: "caramel"},
		{Name: "clementine"},
	}, actualResponse.Items)
	assert.NotEmpty(t, actualResponse.OrderId)
	assert.NotEmpty(t, actualResponse.CreationDate)
	assert.GreaterOrEqual(t, int64(10), responseTimeInMs, "Response time in milliseconds")
}

func TestGetOrder(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	startTime := time.Now()
	apiOrder, err := client.GetApiOrdersOrderId(context.Background(), "IsQah2TkaqS-NONE-DEV-JewgL0Ye73g")
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()
	responseTimeInMs := time.Since(startTime).Milliseconds()

	// THEN
	var actualResponse order_api.OrderResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse order_api.OrderResponse
	readToObject(t, readFile(t, "orderResponse.json"), &expectedResponse)

	assert.Equal(t, 200, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
	assert.GreaterOrEqual(t, int64(10), responseTimeInMs)
}

func TestGetOrderNotFound(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	startTime := time.Now()
	addCorrelationIdHeader := func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Correlation-ID", "2685342d-4888-4d74-9a57-aa5393fc8e35")
		return nil
	}
	apiOrder, err := client.GetApiOrdersOrderId(context.Background(), "NOPE", addCorrelationIdHeader)
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()
	responseTimeInMs := time.Since(startTime).Milliseconds()

	// THEN
	var actualResponse order_api.ErrorResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse order_api.ErrorResponse
	readToObject(t, readFile(t, "orderNotFoundResponse.json"), &expectedResponse)
	expectedResponse.Timestamp = actualResponse.Timestamp

	assert.Equal(t, 404, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
	assert.GreaterOrEqual(t, int64(10), responseTimeInMs, "Response time in milliseconds")
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
