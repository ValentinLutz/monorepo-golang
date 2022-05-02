app.config: config/config.none-dev.yaml ## Copy none-dev config to app directory
	cp config/config.none-dev.yaml app/config.yaml

app.run:: app.config ## Run the app
	cd app && \
		go run main.go

app.build:: ## Build the app into an executable
	cd app && \
		go build
