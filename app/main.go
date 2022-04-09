package main

import (
	"app/api/order"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "golang-reference-project", log.LstdFlags)

	orderAPI := order.NewAPI(logger)

	router := httprouter.New()
	orderAPI.RegisterHandlers(router)

	log.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
