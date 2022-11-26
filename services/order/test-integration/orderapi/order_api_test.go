package orderapi_test

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
	"test-integration/orderapi"
	"testing"
	"time"
)

var config = testingutil.LoadConfig("../../config/test")

func initClient() orderapi.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return orderapi.Client{
		Server: config.BaseURL,
		Client: client,
	}
}

func initDatabase(t *testing.T) *sqlx.DB {
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

func TestGetOrders(t *testing.T) {
	// GIVEN
	client := initClient()
	database := initDatabase(t)

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
	startTime := time.Now()
	apiOrder, err := client.GetApiOrders(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()
	responseTimeInMs := time.Since(startTime).Milliseconds()

	// THEN
	var actualResponse orderapi.OrdersResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.OrdersResponse
	readToObject(t, readFile(t, "ordersResponse.json"), &expectedResponse)

	assert.Equal(t, 200, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
	assert.GreaterOrEqual(t, int64(10), responseTimeInMs, "Response time in milliseconds")
}

func TestPostOrder(t *testing.T) {
	// GIVEN
	client := initClient()
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
	startTime := time.Now()
	apiOrder, err := client.PostApiOrdersWithBody(context.Background(), "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()
	responseTimeInMs := time.Since(startTime).Milliseconds()

	// THEN
	var actualResponse orderapi.OrderResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	assert.Equal(t, 201, apiOrder.StatusCode)
	assert.Equal(t, orderapi.OrderPlaced, actualResponse.Status)
	assert.Equal(t, []orderapi.OrderItemResponse{
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
	database := initDatabase(t)

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
	startTime := time.Now()
	apiOrder, err := client.GetApiOrdersOrderId(context.Background(), "fdCDxjV9o!O-NONE-ZCTH5i6fWcA")
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()
	responseTimeInMs := time.Since(startTime).Milliseconds()

	// THEN
	var actualResponse orderapi.OrderResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.OrderResponse
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
	var actualResponse orderapi.ErrorResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse orderapi.ErrorResponse
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
