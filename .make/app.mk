app/config/config.yaml: ./config/config.none-dev.yaml ## Copy none-dev config to app directory
	install -D config/config.none-dev.yaml app/config/config.yaml

app/serve/openapi/order_api.yaml: api-definition/order_api.yaml ## Copy orders open api definition to app
	install -D api-definition/order_api.yaml app/serve/openapi/order_api.yaml

app/api/order/order.gen.go: api-definition/order_api.yaml api-definition/oapi-codengen.yaml ## Generate orders server from open api definition
	oapi-codegen --config api-definition/oapi-codengen.yaml \
		api-definition/order_api.yaml  > app/api/order/order.gen.go

app.run:: app/config/config.yaml app/serve/openapi/order_api.yaml app/api/order/order.gen.go ## Run the app
	cd app && \
		go run main.go

app.build:: app/serve/openapi/order_api.yaml app/api/order/order.gen.go ## Build the app into an executable
	cd app && \
		go build

app.lint:: app/api/order/order.gen.go ## Runs linters against go code
	cd app && \
		golangci-lint run
