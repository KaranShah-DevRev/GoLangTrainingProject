package main

import (
	controller "CMS/backendServices/deleteService/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/delete-player/{playerId}", func(w http.ResponseWriter, r *http.Request) {
		controller.DeletePlayer()(w, r)
	}).Methods("DELETE")

	router.HandleFunc("/delete-team/{teamId}", func(w http.ResponseWriter, r *http.Request) {
		controller.DeleteTeam()(w, r)
	}).Methods("DELETE")

	router.HandleFunc("/delete-fixture/{fixtureId}", func(w http.ResponseWriter, r *http.Request) {
		controller.DeleteFixture()(w, r)
	}).Methods("DELETE")

	http.ListenAndServe(":8124", router)
}
