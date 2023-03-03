.DEFAULT_GOAL := help

help::
	@egrep -h '\s##\s' $(MAKEFILE_LIST) \
		| awk -F':.*?## | \\| ' '{printf "\033[36m%-38s \033[37m %-68s \033[35m%s \n", $$1, $$2, $$3}'


app/src/incoming/openapi/spec/order_api.yaml: api-definition/order_api.yaml ## Copy order api spec to app
	install -D api-definition/order_api.yaml app/src/incoming/openapi/spec/order_api.yaml


app.run:: app/src/incoming/openapi/spec/order_api.yaml ## Run the app
	cd app && \
		cargo run

app.test:: app/src/incoming/openapi/spec/order_api.yaml ## Test the app
	cd app && \
		cargo test