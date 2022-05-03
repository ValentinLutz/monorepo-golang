app-golang/config.yaml: ./config/config.none-dev.yaml ## Copy none-dev config to app directory
	cp config/config.none-dev.yaml app-golang/config.yaml

app.run:: app-golang/config.yaml ## Run the app
	cd app-golang && \
		go run main.go

app.build:: ## Build the app into an executable
	cd app-golang && \
		go build
