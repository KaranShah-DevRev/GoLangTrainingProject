package main

import (
	"CMS/backendServices/cruService/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	go routes.Routes(router)
	http.ListenAndServe(":8123", router)
}
