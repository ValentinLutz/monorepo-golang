package orderapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"monorepo/libraries/testingutil"
	"monorepo/services/order/test-integration/orderapi"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	testingutil.LoadAndExec(orderapi.GetTestDatabaseInstance(), "truncate_tables.sql")

	code := m.Run()
	os.Exit(code)
}

func Test_GetOrders(t *testing.T) {
	// GIVEN
	client := orderapi.GetOrderApiClientInstance()
	database := orderapi.GetTestDatabaseInstance()

	testingutil.LoadAndExec(database, "init_getOrders.sql")
	customerId, err := uuid.Parse("44bd6239-7e3d-4d4a-90a0-7d4676a00f5c")
	if err != nil {
		t.Fatal(err)
	}

	// WHEN
	apiOrder, err := client.GetOrders(context.Background(), &orderapi.GetOrdersParams{
		CustomerId: &customerId,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(t, 200, apiOrder.StatusCode)

	var actualResponse orderapi.OrdersResponse
	testingutil.ReadToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.OrdersResponse
	testingutil.ReadToObject(t, testingutil.ReadFile(t, "ordersResponse.json"), &expectedResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

func Test_GetOrders_EmptyArray(t *testing.T) {
	// GIVEN
	client := orderapi.GetOrderApiClientInstance()
	customerId, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Fatal(err)
	}

	// WHEN
	apiOrder, err := client.GetOrders(context.Background(), &orderapi.GetOrdersParams{
		CustomerId: &customerId,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(t, 200, apiOrder.StatusCode)

	var actualResponse orderapi.OrdersResponse
	testingutil.ReadToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.OrdersResponse
	testingutil.ReadToObject(t, testingutil.ReadFile(t, "ordersEmptyResponse.json"), &expectedResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

func Test_PostOrder(t *testing.T) {
	// GIVEN
	client := orderapi.GetOrderApiClientInstance()
	orderItems := []orderapi.OrderItemRequest{
		{Name: "caramel"},
		{Name: "clementine"},
	}
	orderRequest := orderapi.OrderRequest{Items: orderItems}
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(orderRequest)
	if err != nil {
		t.Fatal(err)
	}

	// WHEN
	apiOrder, err := client.PostOrdersWithBody(context.Background(), "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(t, 201, apiOrder.StatusCode)

	var actualResponse orderapi.OrderResponse
	testingutil.ReadToObject(t, apiOrder.Body, &actualResponse)
	assert.Equal(t, orderapi.OrderPlaced, actualResponse.Status)
	assert.Equal(t, []orderapi.OrderItemResponse{
		{Name: "caramel"},
		{Name: "clementine"},
	}, actualResponse.Items)
	assert.NotEmpty(t, actualResponse.OrderId)
	assert.NotEmpty(t, actualResponse.CreationDate)
}

func Test_GetOrder(t *testing.T) {
	// GIVEN
	client := orderapi.GetOrderApiClientInstance()
	database := orderapi.GetTestDatabaseInstance()
	testingutil.LoadAndExec(database, "init_getOrder.sql")

	// WHEN
	apiOrder, err := client.GetOrder(context.Background(), "fdCDxjV9o!O-NONE-ZCTH5i6fWcA")
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(t, 200, apiOrder.StatusCode)

	var actualResponse orderapi.OrderResponse
	testingutil.ReadToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.OrderResponse
	testingutil.ReadToObject(t, testingutil.ReadFile(t, "orderResponse.json"), &expectedResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

func Test_GetOrder_NotFound(t *testing.T) {
	// GIVEN
	client := orderapi.GetOrderApiClientInstance()

	// WHEN
	addCorrelationIdHeader := func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Correlation-ID", "2685342d-4888-4d74-9a57-aa5393fc8e35")
		return nil
	}
	apiOrder, err := client.GetOrder(context.Background(), "NOPE", addCorrelationIdHeader)
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(t, 404, apiOrder.StatusCode)

	var actualResponse orderapi.ErrorResponse
	testingutil.ReadToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.ErrorResponse
	testingutil.ReadToObject(t, testingutil.ReadFile(t, "orderNotFoundResponse.json"), &expectedResponse)
	expectedResponse.Timestamp = actualResponse.Timestamp
	assert.Equal(t, expectedResponse, actualResponse)
}
