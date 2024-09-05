package test

import (
	"encoding/json"
	"github.com/google/uuid"
	"monorepo/libraries/testutil"
	testintegration "monorepo/services/order/test-integration"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func Test_Get_Order(t *testing.T) {
	t.Run(
		"get orders for customer", func(t *testing.T) {
			// given
			client := testintegration.GetTestClientInstance()
			database := testintegration.GetTestDatabaseInstance()

			database.MustLoadAndExec("../truncate_tables.sql")
			database.MustLoadAndExec("./init_get_orders.sql")

			// when
			response, err := client.Get("/orders",
				uuid.MustParse("7d5304b9-fac8-4af9-b742-59125ca051f8"),
				map[string]string{
					"customer_id": "44bd6239-7e3d-4d4a-90a0-7d4676a00f5c",
				},
				nil,
			)

			// then
			td.CmpNoError(t, err)
			td.Cmp(t, response.StatusCode, 200)

			body, err := testutil.ParseResponseBodyList(response)
			td.CmpNoError(t, err)

			prettyBody, err := json.MarshalIndent(body, "", "  ")
			td.CmpNoError(t, err)
			t.Log(string(prettyBody))

			td.CmpJSON(t, body, "./get_orders_response.json", []any{})
		},
	)

	t.Run(
		"get orders for customer that has no orders", func(t *testing.T) {
			// given
			client := testintegration.GetTestClientInstance()
			database := testintegration.GetTestDatabaseInstance()

			database.MustLoadAndExec("../truncate_tables.sql")

			// when
			response, err := client.Get("/orders",
				uuid.MustParse("df3c4d1a-48a9-432a-b6c2-6a9a53046993"),
				map[string]string{
					"customer_id": "62a6e56d-e980-4706-96d8-d6cf7b40ad94",
				},
				nil,
			)

			// then
			td.CmpNoError(t, err)
			td.Cmp(t, response.StatusCode, 200)

			body, err := testutil.ParseResponseBodyList(response)
			td.CmpNoError(t, err)

			prettyBody, err := json.MarshalIndent(body, "", "  ")
			td.CmpNoError(t, err)
			t.Log(string(prettyBody))

			td.CmpJSON(t, body, "./get_orders_empty_response.json", []any{})
		},
	)
}
