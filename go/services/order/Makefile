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

dep.oapi-codegen:: # Install oapi-codegen with go install
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4

dep.golangci-lint:: # Install golangci-lint with go install
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.1

dep.kustomize:: # Install kustomize with go install
	go install sigs.k8s.io/kustomize/kustomize/v4@v4.5.6



app/config/config.yaml: config/app/config.none-local.yaml ## Copy none-local config to app directory
	install -D config/app/config.none-local.yaml app/config/config.yaml

app/config/app.key: config/app.key ## Copy none-local private key to app directory
	install -D config/app.key app/config/app.key

app/config/app.crt: config/app.crt ## Copy none-local certificate to app directory
	install -D config/app.crt app/config/app.crt

app/incoming/orderapi/model.gen.go: dep.oapi-codegen api-definition/order_api.yaml api-definition/model.app.yaml ## Generate order models from open api definition
	oapi-codegen --config api-definition/model.app.yaml \
		api-definition/order_api.yaml

app/incoming/orderapi/server.gen.go: dep.oapi-codegen api-definition/order_api.yaml api-definition/server.app.yaml ## Generate order server from open api definition
	oapi-codegen --config api-definition/server.app.yaml \
		api-definition/order_api.yaml

test-integration/orderapi/client.gen.go: dep.oapi-codegen api-definition/order_api.yaml api-definition/client.test.yaml ## Generate integration test orders client from open api definition
	oapi-codegen --config api-definition/client.test.yaml \
		api-definition/order_api.yaml

deployment-kubernetes/base/config.yaml: config/app/config.${PROFILE}.yaml ## Copy config to kubernetes deployment directory
	install -D config/app/config.${PROFILE}.yaml deployment-kubernetes/base/config.yaml



app.run:: app/config/config.yaml app/config/app.key app/config/app.crt app/incoming/orderapi/model.gen.go app/incoming/orderapi/server.gen.go ## Run the app
	cd app && \
		go run main.go

app.build:: app/incoming/orderapi/model.gen.go app/incoming/orderapi/server.gen.go ## Build the app into an executable
	cd app && \
		go build

app.lint:: dep.golangci-lint app/incoming/orderapi/model.gen.go app/incoming/orderapi/server.gen.go ## Runs linters against go code
	cd app && \
		golangci-lint run



test.spectral:: api-definition/order_api.yaml api-definition/ruleset.spectral.yaml
	docker run --rm -it -v ./api-definition/:/tmp stoplight/spectral --ruleset "/tmp/ruleset.spectral.yaml" lint "/tmp/order_api.yaml"

test.unit::  app/incoming/orderapi/model.gen.go app/incoming/orderapi/server.gen.go ## Run the unit tests
	cd app && \
		go test -race -cover ./...

test.smoke:: ## Run the smoke tests
	cd test-smoke && \
		go test -count=1 ./...

test.integration:: test-integration/orderapi/client.gen.go ## Run the integration tests
	cd test-integration && \
		go test -count=1 ./...

test.load:: ## Run load tests
	k6 run test-load/script.js 



database.migrate:: ## Migrate database | PROFILE, FLYWAY_USER, FLYWAY_PASSWORD
	cd migration-database && \
		flyway clean \
		migrate \
		-configFiles=${PROFILE}.properties \
		-user=${FLYWAY_USER} \
		-password=${FLYWAY_PASSWORD}



docker.build:: app/incoming/orderapi/server.gen.go ## Build container image | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
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



kubernetes.image:: dep.kustomize ## Update app image version
	cd deployment-kubernetes/base && \
		kustomize edit set image ghcr.io/valentinlutz/order-service=ghcr.io/valentinlutz/order-service:${VERSION}

kubernetes.build:: dep.kustomize deployment-kubernetes/base/config.yaml ## Build kubernetes deployment
	kustomize build deployment-kubernetes/base > deployment-kubernetes/${PROFILE}.yaml

kubernetes.deploy:: kubernetes.build ## Deploy app to kubernetes
	kubectl apply -f deployment-kubernetes/${PROFILE}.yaml