app/config/config.yaml: ./config/config.none-dev.yaml ## Copy none-dev config to app directory
	install -D config/config.none-dev.yaml app/config/config.yaml

app/serve/openapi/orders.yaml: api-definition/orders.yaml ## Copy orders open api definition to app
	install -D api-definition/orders.yaml app/serve/openapi/orders.yaml

app/api/orders/orders.gen.go: api-definition/orders.yaml ## Generate orders server from open api definition
	oapi-codegen -generate types \
		-package orders \
		./api-definition/orders.yaml  > app/api/orders/orders.gen.go

app.run:: app/config/config.yaml app/serve/openapi/orders.yaml app/api/orders/orders.gen.go ## Run the app
	cd app && \
		go run main.go

app.build:: app/serve/openapi/orders.yaml app/api/orders/orders.gen.go ## Build the app into an executable
	cd app && \
		go build

app.lint:: app/api/orders/orders.gen.go ## Runs linters against go code
	cd app && \
		golangci-lint run
