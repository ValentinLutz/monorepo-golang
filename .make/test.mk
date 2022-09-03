test-integration/order_api/order.gen.go: api-definition/order_api.yaml api-definition/oapi-codengen.yaml ## Generate integration test orders client from open api definition
	oapi-codegen --config api-definition/oapi-codengen-test.yaml \
		api-definition/order_api.yaml  > test-integration/order_api/order.gen.go

test.unit::  app/serve/openapi/order_api.yaml app/adapter/order_api/order.gen.go ## Run the unit tests
	cd app && \
		go test ./...

test.smoke:: ## Run the smoke tests
	cd test-smoke && \
		go test -count=1 ./...

test.integration:: test-integration/order_api/order.gen.go ## Run the integration tests
	cd test-integration && \
		go test -count=1 ./...

test.load:: ## Run load tests
	docker run -it \
		--rm \
		--volume ${PWD}/test-load:/k6 \
		--network golang-reference-project \
        grafana/k6:0.39.0 \
		run /k6/script.js \

test:: test.unit test.smoke test.integration test.load ## Run all tests