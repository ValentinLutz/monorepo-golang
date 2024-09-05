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
		"create new order", func(t *testing.T) {
			// given
			client := testintegration.GetTestClientInstance()
			database := testintegration.GetTestDatabaseInstance()

			database.MustLoadAndExec("../truncate_tables.sql")
			requestFile := testutil.ReadFile(t, "./post_order_request.json")

			// when
			response, err := client.Post("/orders",
				uuid.MustParse("572694b9-549a-44a9-9477-1ae21cbda887"),
				nil,
				requestFile,
			)

			// then
			td.CmpNoError(t, err)
			td.Cmp(t, response.StatusCode, 201)

			body, err := testutil.ParseResponseBody(response)
			td.CmpNoError(t, err)

			prettyBody, err := json.MarshalIndent(body, "", "  ")
			td.CmpNoError(t, err)
			t.Log(string(prettyBody))

			td.CmpJSON(t, body, "./post_order_response.json", []any{})
		},
	)
}
