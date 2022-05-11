test-integration/orders/orders.gen.go: api-definition/orders.yaml ## Generate integration test orders client from open api definition
	oapi-codegen -generate types,client \
		-package orders \
		./api-definition/orders.yaml  > test-integration/orders/orders.gen.go

test.unit::  ## Run the unit tests
	cd app && \
		go test ./...

test.smoke:: ## Run the smoke tests
	cd test-smoke && \
		go test -count=1 ./...

test.integration:: test-integration/orders/orders.gen.go ## Run the integration tests
	cd test-integration && \
		go test -count=1 ./...
