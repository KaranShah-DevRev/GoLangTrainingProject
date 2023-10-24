package routes

import (
	fixturecontroller "CMS/backendServices/cruService/controllers/fixtureController"
	playerController "CMS/backendServices/cruService/controllers/playerController"
	teamController "CMS/backendServices/cruService/controllers/teamController"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes(router *mux.Router) {

	// Player routes
	router.HandleFunc("/create-player", playerController.CreatePlayerHandler()).Methods("POST")
	router.HandleFunc("/get-player/{playerId}", func(w http.ResponseWriter, r *http.Request) {
		playerController.GetPlayerHandler()(w, r)
	}).Methods("GET")
	router.HandleFunc("/update-player/{playerId}", func(w http.ResponseWriter, r *http.Request) {
		playerController.UpdatePlayerHandler()(w, r)
	}).Methods("PUT")
	router.HandleFunc("/get-players", playerController.GetAllPlayers()).Methods("GET")

	// Team routes
	router.HandleFunc("/create-team", teamController.CreateTeamHandler()).Methods("POST")
	router.HandleFunc("/get-team/{teamId}", func(w http.ResponseWriter, r *http.Request) {
		teamController.GetTeamHandler()(w, r)
	}).Methods("GET")
	router.HandleFunc("/update-team/{teamId}", func(w http.ResponseWriter, r *http.Request) {
		teamController.UpdateTeamHandler()(w, r)
	}).Methods("PUT")
	router.HandleFunc("/get-teams", teamController.GetAllTeamsHandler()).Methods("GET")

	// Fixture routes
	router.HandleFunc("/create-fixture", fixturecontroller.CreateFixtureHandler()).Methods("POST")
	router.HandleFunc("/get-fixture/{fixtureId}", func(w http.ResponseWriter, r *http.Request) {
		fixturecontroller.GetFixtureHandler()(w, r)
	}).Methods("GET")
	router.HandleFunc("/update-fixture/{fixtureId}", func(w http.ResponseWriter, r *http.Request) {
		fixturecontroller.UpdateFixtureHandler()(w, r)
	}).Methods("PUT")
	router.HandleFunc("/get-fixtures", fixturecontroller.GetAllFixturesHandler()).Methods("GET")
}
