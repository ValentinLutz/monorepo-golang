docker.build:: app/serve/openapi/orders.yaml app/api/orders/orders.gen.go ## Build container images | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
	docker build \
    		-t ${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:${VERSION} \
    		-t ${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:latest \
    		app

docker.push:: ## Publish container images | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
	docker push \
		${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:${VERSION}
	docker push \
		${DOCKER_REGISTRY}/${DOCKER_REPOSITORY}/${PROJECT_NAME}:latest

docker.up:: ## Start containers | PROJECT_NAME
	docker-compose -p ${PROJECT_NAME} \
		-f deployment-docker/docker-compose.yaml \
		up -d

docker.down:: ## Shutdown containers | PROJECT_NAME
	docker-compose -p ${PROJECT_NAME} \
		-f deployment-docker/docker-compose.yaml \
		down

docker.app.up:: docker.build ## Start app container | PROJECT_NAME
	docker-compose -p ${PROJECT_NAME} \
		-f deployment-docker/app/docker-compose.yaml \
		up