package testintegration

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"monorepo/libraries/testingutil"
	"monorepo/services/order/test-integration/orderapi"
	"net/http"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func newOrderApiClient(t *testing.T, config *testingutil.Config) *orderapi.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	basicAuth, err := securityprovider.NewSecurityProviderBasicAuth("test", "test")
	if err != nil {
		t.Fatal(err)
	}

	orderApiClient, err := orderapi.NewClient(
		config.BaseURL+"/",
		orderapi.WithHTTPClient(client),
		orderapi.WithRequestEditorFn(basicAuth.Intercept),
	)
	if err != nil {
		t.Fatal(err)
	}

	return orderApiClient
}

func (suite *IntegrationTestSuite) Test_GetOrders() {
	// GIVEN
	client := newOrderApiClient(suite.T(), suite.config)
	testingutil.LoadAndExec(suite.T(), suite.database, "orderapi/init_getOrders.sql")
	customerId, err := uuid.Parse("44bd6239-7e3d-4d4a-90a0-7d4676a00f5c")
	if err != nil {
		suite.T().Fatal(err)
	}

	// WHEN
	apiOrder, err := client.GetOrders(context.Background(), &orderapi.GetOrdersParams{
		CustomerId: &customerId,
	})
	if err != nil {
		suite.T().Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(suite.T(), 200, apiOrder.StatusCode)

	var actualResponse orderapi.OrdersResponse
	testingutil.ReadToObject(suite.T(), apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.OrdersResponse
	testingutil.ReadToObject(suite.T(), testingutil.ReadFile(suite.T(), "orderapi/ordersResponse.json"), &expectedResponse)
	assert.Equal(suite.T(), expectedResponse, actualResponse)
}

func (suite *IntegrationTestSuite) Test_GetOrders_EmptyArray() {
	// GIVEN
	client := newOrderApiClient(suite.T(), suite.config)
	customerId, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if err != nil {
		suite.T().Fatal(err)
	}

	// WHEN
	apiOrder, err := client.GetOrders(context.Background(), &orderapi.GetOrdersParams{
		CustomerId: &customerId,
	})
	if err != nil {
		suite.T().Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(suite.T(), 200, apiOrder.StatusCode)

	var actualResponse orderapi.OrdersResponse
	testingutil.ReadToObject(suite.T(), apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.OrdersResponse
	testingutil.ReadToObject(suite.T(), testingutil.ReadFile(suite.T(), "orderapi/ordersEmptyResponse.json"), &expectedResponse)
	assert.Equal(suite.T(), expectedResponse, actualResponse)
}

func (suite *IntegrationTestSuite) Test_PostOrder() {
	// GIVEN
	client := newOrderApiClient(suite.T(), suite.config)
	orderItems := []orderapi.OrderItemRequest{
		{Name: "caramel"},
		{Name: "clementine"},
	}
	orderRequest := orderapi.OrderRequest{Items: orderItems}
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(orderRequest)
	if err != nil {
		suite.T().Fatal(err)
	}

	// WHEN
	apiOrder, err := client.PostOrdersWithBody(context.Background(), "application/json", &body)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(suite.T(), 201, apiOrder.StatusCode)

	var actualResponse orderapi.OrderResponse
	testingutil.ReadToObject(suite.T(), apiOrder.Body, &actualResponse)
	assert.Equal(suite.T(), orderapi.OrderPlaced, actualResponse.Status)
	assert.Equal(suite.T(), []orderapi.OrderItemResponse{
		{Name: "caramel"},
		{Name: "clementine"},
	}, actualResponse.Items)
	assert.NotEmpty(suite.T(), actualResponse.OrderId)
	assert.NotEmpty(suite.T(), actualResponse.CreationDate)
}

func (suite *IntegrationTestSuite) Test_GetOrder() {
	// GIVEN
	client := newOrderApiClient(suite.T(), suite.config)
	testingutil.LoadAndExec(suite.T(), suite.database, "orderapi/init_getOrder.sql")

	// WHEN
	apiOrder, err := client.GetOrder(context.Background(), "fdCDxjV9o!O-NONE-ZCTH5i6fWcA")
	if err != nil {
		suite.T().Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(suite.T(), 200, apiOrder.StatusCode)

	var actualResponse orderapi.OrderResponse
	testingutil.ReadToObject(suite.T(), apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.OrderResponse
	testingutil.ReadToObject(suite.T(), testingutil.ReadFile(suite.T(), "orderapi/orderResponse.json"), &expectedResponse)
	assert.Equal(suite.T(), expectedResponse, actualResponse)
}

func (suite *IntegrationTestSuite) Test_GetOrder_NotFound() {
	// GIVEN
	client := newOrderApiClient(suite.T(), suite.config)

	// WHEN
	addCorrelationIdHeader := func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Correlation-ID", "2685342d-4888-4d74-9a57-aa5393fc8e35")
		return nil
	}
	apiOrder, err := client.GetOrder(context.Background(), "NOPE", addCorrelationIdHeader)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(suite.T(), 404, apiOrder.StatusCode)

	var actualResponse orderapi.ErrorResponse
	testingutil.ReadToObject(suite.T(), apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.ErrorResponse
	testingutil.ReadToObject(suite.T(), testingutil.ReadFile(suite.T(), "orderapi/orderNotFoundResponse.json"), &expectedResponse)
	expectedResponse.Timestamp = actualResponse.Timestamp
	assert.Equal(suite.T(), expectedResponse, actualResponse)
}
