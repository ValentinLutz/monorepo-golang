docker.build:: app/serve/openapi/order_api.yaml app/adapter/order_api/order.gen.go ## Build container images | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
	docker build \
    		-t ${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:${VERSION} \
    		-t ${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:latest \
    		app

docker.push:: ## Publish container images | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
	docker push \
		${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:${VERSION}
	docker push \
		${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:latest

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