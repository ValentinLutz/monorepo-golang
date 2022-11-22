package order_api_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/ValentinLutz/monrepo/libraries/testingutil"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"test-integration/order_api"
	"testing"
	"time"
)

var config = testingutil.LoadConfig("../../config/test")

const addOrders = `
INSERT INTO order_service.order
(order_id, workflow, creation_date, order_status)
VALUES ('IsQah2TkaqS-NONE-DEV-JewgL0Ye73g', 'default_workflow', '1980-01-01 00:00:00 +00:00', 'order_placed'),

       ('Fs2VoM7ZhrK-NONE-DEV-vzTf7kaHbRA', 'default_workflow', '1980-01-01 00:00:00 +00:00', 'order_in_progress'),

       ('sgy1K3*SXcv-NONE-DEV-eVbldUAYXnA', 'default_workflow', '1980-01-01 00:00:00 +00:00', 'order_canceled'),

       ('F2P!criGu2L-NONE-DEV-fJ7bBFx1vHg', 'default_workflow', '1980-01-01 00:00:00 +00:00', 'order_completed');

INSERT INTO order_service.order_item
    (order_id, creation_date, item_name)
VALUES ('IsQah2TkaqS-NONE-DEV-JewgL0Ye73g', '1980-01-01 00:00:00 +00:00', 'orange'),
       ('IsQah2TkaqS-NONE-DEV-JewgL0Ye73g', '1980-01-01 00:00:00 +00:00', 'banana'),

       ('Fs2VoM7ZhrK-NONE-DEV-vzTf7kaHbRA', '1980-01-01 00:00:00 +00:00', 'chocolate'),

       ('sgy1K3*SXcv-NONE-DEV-eVbldUAYXnA', '1980-01-01 00:00:00 +00:00', 'marshmallow'),

       ('F2P!criGu2L-NONE-DEV-fJ7bBFx1vHg', '1980-01-01 00:00:00 +00:00', 'apple');

`

const addOrder = `
INSERT INTO order_service.order
(order_id, workflow, creation_date, order_status)
VALUES ('fdCDxjV9o!O-NONE-DEV-ZCTH5i6fWcA', 'default_workflow', '1980-01-01 00:00:00 +00:00', 'order_placed');

INSERT INTO order_service.order_item
    (order_id, creation_date, item_name)
VALUES ('fdCDxjV9o!O-NONE-DEV-ZCTH5i6fWcA', '1980-01-01 00:00:00 +00:00', 'orange'),
       ('fdCDxjV9o!O-NONE-DEV-ZCTH5i6fWcA', '1980-01-01 00:00:00 +00:00', 'banana');
`

const cleanDatabase = `
DELETE FROM order_service.order_item;
DELETE FROM order_service.order;
`

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

func initDatabase() *sqlx.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.Username, config.Database.Password, config.Database.Database,
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(cleanDatabase)
	if err != nil {
		panic(err)
	}
	return db
}

func createOrders(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec(addOrders)
	if err != nil {
		t.Fatal(err)
	}
}

func createOrder(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec(addOrder)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOrders(t *testing.T) {
	// GIVEN
	client := initClient()
	database := initDatabase()
	createOrders(t, database)

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
	assert.GreaterOrEqual(t, int64(50), responseTimeInMs, "Response time in milliseconds")
}

func TestGetOrder(t *testing.T) {
	// GIVEN
	client := initClient()
	database := initDatabase()
	createOrder(t, database)

	// WHEN
	startTime := time.Now()
	apiOrder, err := client.GetApiOrdersOrderId(context.Background(), "fdCDxjV9o!O-NONE-DEV-ZCTH5i6fWcA")
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
	assert.GreaterOrEqual(t, int64(50), responseTimeInMs)
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
