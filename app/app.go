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

	productsHandler := api.NewOrderApi(logger)

	router := httprouter.New()
	productsHandler.RegisterHandlers(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
