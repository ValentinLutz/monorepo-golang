docker.build:: app/serve/openapi/order_api.yaml app/adapter/order_api/order.gen.go ## Build container image | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
ifneq ($(findstring SNAPSHOT,$(VERSION)),SNAPSHOT)
	docker build \
		-t ${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:latest \
		app
endif
	docker build \
		-t ${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:${VERSION} \
		app

docker.push:: ## Publish container image | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
ifneq ($(findstring SNAPSHOT,$(VERSION)),SNAPSHOT)
	docker push \
		${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:latest
endif
	docker push \
		${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:${VERSION}

docker.up:: docker.build ## Start containers | PROJECT_NAME
	docker-compose -p ${PROJECT_NAME} \
		-f deployment-docker/docker-compose.yaml \
		up -d --force-recreate

docker.down:: ## Shutdown containers | PROJECT_NAME
	docker-compose -p ${PROJECT_NAME} \
		-f deployment-docker/docker-compose.yaml \
		down

docker.app.up:: docker.up database.migrate ## Start containers | PROJECT_NAME
	docker logs app \
		--follow