docker.build:: ## Build the app docker image | PROJECT_NAME
	docker build \
		-t ${PROJECT_NAME}:latest \
		app-golang

docker.up:: ## Start docker containers | PROJECT_NAME
	docker-compose -p ${PROJECT_NAME} \
		-f deployment-docker/docker-compose.yaml \
		up -d

docker.down:: ## Shutdown docker containers | PROJECT_NAME
	docker-compose -p ${PROJECT_NAME} \
		-f deployment-docker/docker-compose.yaml \
		down

docker.app.up:: docker.build ## Start app docker container | PROJECT_NAME
	docker-compose -p ${PROJECT_NAME} \
		-f deployment-docker/app/docker-compose.yaml \
		up