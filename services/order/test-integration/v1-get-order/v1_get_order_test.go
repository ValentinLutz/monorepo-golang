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
		"get order", func(t *testing.T) {
			// given
			client := testintegration.GetTestClientInstance()
			database := testintegration.GetTestDatabaseInstance()

			database.MustLoadAndExec("../truncate_tables.sql")
			database.MustLoadAndExec("./init_get_order.sql")

			// when
			response, err := client.Get("/orders/01J71WT3T81XK-NONE-X7Y3DD7FGWV0F", uuid.MustParse("7d5304b9-fac8-4af9-b742-59125ca051f8"), nil, nil)

			// then
			td.CmpNoError(t, err)
			td.Cmp(t, response.StatusCode, 200)

			body, err := testutil.ParseResponseBody(response)
			td.CmpNoError(t, err)

			prettyBody, err := json.MarshalIndent(body, "", "  ")
			td.CmpNoError(t, err)
			t.Log(string(prettyBody))

			td.CmpJSON(t, body, "./get_order_response.json", []any{})
		},
	)

	t.Run(
		"get order that does not exists", func(t *testing.T) {
			// given
			client := testintegration.GetTestClientInstance()
			database := testintegration.GetTestDatabaseInstance()

			database.MustLoadAndExec("../truncate_tables.sql")

			// when
			response, err := client.Get("/orders/01J71WVCVV6PZ-NONE-84HWHJDBB8K72", uuid.MustParse("df3c4d1a-48a9-432a-b6c2-6a9a53046993"), nil, nil)

			// then
			td.CmpNoError(t, err)
			td.Cmp(t, response.StatusCode, 404)

			body, err := testutil.ParseResponseBody(response)
			td.CmpNoError(t, err)

			prettyBody, err := json.MarshalIndent(body, "", "  ")
			td.CmpNoError(t, err)
			t.Log(string(prettyBody))

			td.CmpJSON(t, body, "./get_order_not_found_response.json", []any{})
		},
	)
}
