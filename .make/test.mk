test-integration/orders/orders.gen.go: api-definition/orders.yaml api-definition/oapi-codengen.yaml ## Generate integration test orders client from open api definition
	oapi-codegen --config api-definition/oapi-codengen-test.yaml \
		api-definition/orders.yaml  > test-integration/order/orders.gen.go

test.unit::  app/serve/openapi/orders.yaml app/api/orders/orders.gen.go ## Run the unit tests
	cd app && \
		go test ./...

test.smoke:: ## Run the smoke tests
	cd test-smoke && \
		go test -count=1 ./...

test.integration:: test-integration/orders/orders.gen.go ## Run the integration tests
	cd test-integration && \
		go test -count=1 ./...

test:: test.unit test.smoke test.integration ## Run all tests