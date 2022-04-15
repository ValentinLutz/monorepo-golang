app.run:: ## Run the app
	cd app && \
		go run main.go

app.build:: ## Build the app into an executable
	cd app && \
		go build
