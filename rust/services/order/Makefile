.DEFAULT_GOAL := help

export DATABASE_URL ?= postgres://test:test@localhost:5432/test
export PROFILE ?= none-local
export FLYWAY_USER ?= test
export FLYWAY_PASSWORD ?= test


help::
	@egrep -h '\s##\s' $(MAKEFILE_LIST) \
		| awk -F':.*?## | \\| ' '{printf "\033[36m%-38s \033[37m %-68s \033[35m%s \n", $$1, $$2, $$3}'


app/src/incoming/openapi/spec/order_api.yaml: api-definition/order_api.yaml ## Copy order api spec to app
	install -D api-definition/order_api.yaml app/src/incoming/openapi/spec/order_api.yaml


app.run:: app/src/incoming/openapi/spec/order_api.yaml ## Run the app
	cd app && \
		cargo run


test.unit:: app/src/incoming/openapi/spec/order_api.yaml ## Test the app
	cd app && \
		cargo test

test.load:: ## Run load tests
	k6 run test-load/script.js


database.migrate:: ## Migrate database | PFOFILE, FLYWAY_USER, FLYWAY_PASSWORD
	cd migration-database && \
		flyway clean \
		migrate \
		-configFiles=${PROFILE}.properties \
		-user=${FLYWAY_USER} \
		-password=${FLYWAY_PASSWORD}



docker.up:: ## Start containers
	docker compose \
		-f deployment-docker/docker-compose.yaml \
		--force-recreate