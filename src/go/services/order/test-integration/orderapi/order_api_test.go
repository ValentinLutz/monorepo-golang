package orderapi_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"io"
	"monorepo/libraries/testingutil"
	"monorepo/services/order/test-integration/orderapi"
	"net/http"
	"os"
	"testing"
)

var config = testingutil.LoadConfig("../../config/test")

func newClient(t *testing.T) *orderapi.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	basicAuth, err := securityprovider.NewSecurityProviderBasicAuth("test", "test")
	if err != nil {
		t.Fatal(err)
	}

	orderApiClient, err := orderapi.NewClient(
		config.BaseURL+"/api",
		orderapi.WithHTTPClient(client),
		orderapi.WithRequestEditorFn(basicAuth.Intercept),
	)
	if err != nil {
		t.Fatal(err)
	}

	return orderApiClient
}

func newDatabase(t *testing.T) *sqlx.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.Username, config.Database.Password, config.Database.Database,
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		t.Fatal(err)
	}

	cleanDatabase := `
TRUNCATE TABLE order_service.order_item, order_service.order;
`
	_, err = db.Exec(cleanDatabase)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func exec(t *testing.T, db *sqlx.DB, query string) {
	_, err := db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetOrders(t *testing.T) {
	// GIVEN
	client := newClient(t)
	database := newDatabase(t)

	addOrders := `
INSERT INTO order_service.order
(order_id, workflow, creation_date, order_status)
VALUES ('IsQah2TkaqS-NONE-JewgL0Ye73g', 'default_workflow', '1980-01-01 00:00:00 +00:00', 'order_placed'),

	   ('Fs2VoM7ZhrK-NONE-vzTf7kaHbRA', 'default_workflow', '1980-01-01 00:00:00 +00:00', 'order_in_progress'),

	   ('sgy1K3*SXcv-NONE-eVbldUAYXnA', 'default_workflow', '1980-01-01 00:00:00 +00:00', 'order_canceled'),

	   ('F2P!criGu2L-NONE-fJ7bBFx1vHg', 'default_workflow', '1980-01-01 00:00:00 +00:00', 'order_completed');

INSERT INTO order_service.order_item
	(order_id, creation_date, item_name)
VALUES ('IsQah2TkaqS-NONE-JewgL0Ye73g', '1980-01-01 00:00:00 +00:00', 'orange'),
	   ('IsQah2TkaqS-NONE-JewgL0Ye73g', '1980-01-01 00:00:00 +00:00', 'banana'),

	   ('Fs2VoM7ZhrK-NONE-vzTf7kaHbRA', '1980-01-01 00:00:00 +00:00', 'chocolate'),

	   ('sgy1K3*SXcv-NONE-eVbldUAYXnA', '1980-01-01 00:00:00 +00:00', 'marshmallow'),

	   ('F2P!criGu2L-NONE-fJ7bBFx1vHg', '1980-01-01 00:00:00 +00:00', 'apple');
`
	exec(t, database, addOrders)

	// WHEN
	apiOrder, err := client.GetOrders(context.Background(), &orderapi.GetOrdersParams{})
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(t, 200, apiOrder.StatusCode)

	var actualResponse orderapi.OrdersResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.OrdersResponse
	readToObject(t, readFile(t, "ordersResponse.json"), &expectedResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

func Test_GetOrders_EmptyArray(t *testing.T) {
	// GIVEN
	client := newClient(t)
	_ = newDatabase(t)

	// WHEN
	apiOrder, err := client.GetOrders(context.Background(), &orderapi.GetOrdersParams{})
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(t, 200, apiOrder.StatusCode)

	var actualResponse orderapi.OrdersResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.OrdersResponse
	readToObject(t, readFile(t, "ordersEmptyResponse.json"), &expectedResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

func Test_PostOrder(t *testing.T) {
	// GIVEN
	client := newClient(t)
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
	readToObject(t, apiOrder.Body, &actualResponse)
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
	client := newClient(t)
	database := newDatabase(t)

	const addOrder = `
INSERT INTO order_service.order
(order_id, workflow, creation_date, order_status)
VALUES ('fdCDxjV9o!O-NONE-ZCTH5i6fWcA', 'default_workflow', '1980-01-01 00:00:00 +00:00', 'order_placed');

INSERT INTO order_service.order_item
    (order_id, creation_date, item_name)
VALUES ('fdCDxjV9o!O-NONE-ZCTH5i6fWcA', '1980-01-01 00:00:00 +00:00', 'orange'),
       ('fdCDxjV9o!O-NONE-ZCTH5i6fWcA', '1980-01-01 00:00:00 +00:00', 'banana');
`

	exec(t, database, addOrder)

	// WHEN
	apiOrder, err := client.GetOrder(context.Background(), "fdCDxjV9o!O-NONE-ZCTH5i6fWcA")
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	assert.Equal(t, 200, apiOrder.StatusCode)

	var actualResponse orderapi.OrderResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.OrderResponse
	readToObject(t, readFile(t, "orderResponse.json"), &expectedResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

func Test_GetOrder_NotFound(t *testing.T) {
	// GIVEN
	client := newClient(t)

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
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.ErrorResponse
	readToObject(t, readFile(t, "orderNotFoundResponse.json"), &expectedResponse)
	expectedResponse.Timestamp = actualResponse.Timestamp
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
