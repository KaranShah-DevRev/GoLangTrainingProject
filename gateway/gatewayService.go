package main

import (
	"CMS/gateway/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./gateway.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := ":" + os.Getenv("PORT")
	router := mux.NewRouter()
	go routes.Routes(router)
	fmt.Print("Gateway running on port " + port + "\n")
	http.ListenAndServe(port, router)
}
