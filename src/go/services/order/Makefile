.DEFAULT_GOAL := help

help::
	@egrep -h '\s##\s' $(MAKEFILE_LIST) \
		| awk -F':.*?## | \\| ' '{printf "\033[36m%-38s \033[37m %-68s \033[35m%s \n", $$1, $$2, $$3}'

export PROJECT_NAME ?= order-service
export VERSION ?= latest
export PROFILE ?= none-local
export FLYWAY_USER ?= test
export FLYWAY_PASSWORD ?= test
export DOCKER_REGISTRY ?= ghcr.io
export DOCKER_REPOSITORY ?= valentinlutz



app/config/config.yaml: config/app/config.none-local.yaml ## Copy none-local config to app directory
	install -D config/app/config.none-local.yaml app/config/config.yaml

app/config/app.key: config/app.key ## Copy none-local private key to app directory
	install -D config/app.key app/config/app.key

app/config/app.crt: config/app.crt ## Copy none-local certificate to app directory
	install -D config/app.crt app/config/app.crt

app/adapter/orderapi/order.gen.go: api-definition/order_api.yaml api-definition/oapi-codegen-app.yaml ## Generate orders server from open api definition
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.1 && \
		oapi-codegen --config api-definition/oapi-codegen-app.yaml \
		api-definition/order_api.yaml  > app/adapter/orderapi/order.gen.go

app.run:: app/config/config.yaml app/config/app.key app/config/app.crt app/adapter/orderapi/order.gen.go ## Run the app
	cd app && \
		go run main.go

app.build:: app/adapter/orderapi/order.gen.go ## Build the app into an executable
	cd app && \
		go build

app.lint:: app/adapter/orderapi/order.gen.go ## Runs linters against go code
	cd app && \
		golangci-lint run



test-integration/orderapi/order.gen.go: api-definition/order_api.yaml api-definition/oapi-codegen-test.yaml ## Generate integration test orders client from open api definition
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.1 && \
		oapi-codegen --config api-definition/oapi-codegen-test.yaml \
		api-definition/order_api.yaml  > test-integration/orderapi/order.gen.go

test.unit::  app/adapter/orderapi/order.gen.go ## Run the unit tests
	cd app && \
		go test -cover ./...

test.smoke:: ## Run the smoke tests
	cd test-smoke && \
		go test -count=1 ./...

test.integration:: test-integration/orderapi/order.gen.go ## Run the integration tests
	cd test-integration && \
		go test -count=1 ./...

test.load:: ## Run load tests
	docker run -it \
		--rm \
		--volume ${PWD}/test-load:/k6 \
		--network order-service \
        grafana/k6:0.39.0 \
		run /k6/script.js \

test:: test.unit test.smoke test.integration test.load ## Run all tests



database.migrate:: ## Migrate database | PROFILE, FLYWAY_USER, FLYWAY_PASSWORD
	cd migration-database && \
		flyway clean \
		migrate \
		-configFiles=${PROFILE}.properties \
		-user=${FLYWAY_USER} \
		-password=${FLYWAY_PASSWORD}



docker.build:: app/adapter/orderapi/order.gen.go ## Build container image | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
ifneq ($(findstring SNAPSHOT,$(VERSION)),SNAPSHOT)
	docker build \
		-f ./app/Dockerfile \
		-t ${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:latest \
		../../
endif
	docker build \
		-f ./app/Dockerfile \
		-t ${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:${VERSION} \
		../../

docker.push:: ## Publish container image | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
ifneq ($(findstring SNAPSHOT,$(VERSION)),SNAPSHOT)
	docker push \
		${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:latest
endif
	docker push \
		${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:${VERSION}

docker.up:: ## Start containers | PROJECT_NAME
	docker compose -p ${PROJECT_NAME} \
		-f deployment-docker/docker-compose.yaml \
		up -d --force-recreate

docker.down:: ## Shutdown containers | PROJECT_NAME
	docker compose -p ${PROJECT_NAME} \
		-f deployment-docker/docker-compose.yaml \
		down

docker.app.up:: docker.up docker.build database.migrate ## Start containers | PROJECT_NAME
	docker compose -p ${PROJECT_NAME} \
    		-f deployment-docker/docker-compose-app.yaml \
    		up -d --force-recreate
	docker logs app \
		--follow