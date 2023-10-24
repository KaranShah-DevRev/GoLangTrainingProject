package routes

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func forwardRequest(url string) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Accept", "application/json")

		switch request.Method {
		case "GET":
			resp, err := http.Get(url)
			if err != nil {
				log.Println("Error in request")
				return
			}
			defer resp.Body.Close()
			for header, values := range resp.Header {
				response.Header()[header] = values
			}
			response.WriteHeader(resp.StatusCode)
			_, _ = io.Copy(response, resp.Body)
		case "POST":
			resp, err := http.Post(url, "application/json", request.Body)
			if err != nil {
				log.Println("Error in request")
				return
			}
			defer resp.Body.Close()
			for header, values := range resp.Header {
				response.Header()[header] = values
			}
			response.WriteHeader(resp.StatusCode)
			_, _ = io.Copy(response, resp.Body)
		case "PUT":
			req, err := http.NewRequest("PUT", url, request.Body)
			if err != nil {
				log.Println("Error in request")
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("Error in request")
				return
			}
			defer resp.Body.Close()
			for header, values := range resp.Header {
				response.Header()[header] = values
			}
			response.WriteHeader(resp.StatusCode)
			_, _ = io.Copy(response, resp.Body)
		case "DELETE":
			req, err := http.NewRequest("DELETE", url, request.Body)
			if err != nil {
				log.Println("Error in request")
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("Error in request")
				return
			}
			defer resp.Body.Close()
			for header, values := range resp.Header {
				response.Header()[header] = values
			}
			response.WriteHeader(resp.StatusCode)
			_, _ = io.Copy(response, resp.Body)
		}
	}
}

func Routes(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Cricket Management System!"))
	}).Methods("GET")

	//Create, Read, Update, Delete Player
	router.HandleFunc("/create-player", forwardRequest("http://localhost:8123/create-player")).Methods("POST")
	router.HandleFunc("/get-player/{playerId}", func(w http.ResponseWriter, r *http.Request) {
		url := "http://localhost:8123/get-player/" + mux.Vars(r)["playerId"]
		forwardRequest(url)(w, r)
	}).Methods("GET")
	router.HandleFunc("/update-player/{playerId}", func(w http.ResponseWriter, r *http.Request) {
		url := "http://localhost:8123/update-player/" + mux.Vars(r)["playerId"]
		forwardRequest(url)(w, r)
	}).Methods("PUT")
	router.HandleFunc("/delete-player/{playerId}", func(w http.ResponseWriter, r *http.Request) {
		url := "http://localhost:8124/delete-player/" + mux.Vars(r)["playerId"]
		forwardRequest(url)(w, r)
	}).Methods("DELETE")
	router.HandleFunc("/get-players", forwardRequest("http://localhost:8123/get-players")).Methods("GET")

	//Create, Read, Update, Delete Team
	router.HandleFunc("/create-team", forwardRequest("http://localhost:8123/create-team")).Methods("POST")
	router.HandleFunc("/get-team/{teamId}", func(w http.ResponseWriter, r *http.Request) {
		url := "http://localhost:8123/get-team/" + mux.Vars(r)["teamId"]
		forwardRequest(url)(w, r)
	}).Methods("GET")
	router.HandleFunc("/update-team/{teamId}", func(w http.ResponseWriter, r *http.Request) {
		url := "http://localhost:8123/update-team/" + mux.Vars(r)["teamId"]
		forwardRequest(url)(w, r)
	}).Methods("PUT")
	router.HandleFunc("/delete-team/{teamId}", func(w http.ResponseWriter, r *http.Request) {
		url := "http://localhost:8124/delete-team/" + mux.Vars(r)["teamId"]
		forwardRequest(url)(w, r)
	}).Methods("DELETE")
	router.HandleFunc("/get-teams", forwardRequest("http://localhost:8123/get-teams")).Methods("GET")

	//Create, Read, Update, Delete Fixture
	router.HandleFunc("/create-fixture", forwardRequest("http://localhost:8123/create-fixture")).Methods("POST")
	router.HandleFunc("/get-fixture/{fixtureId}", func(w http.ResponseWriter, r *http.Request) {
		url := "http://localhost:8123/get-fixture/" + mux.Vars(r)["fixtureId"]
		forwardRequest(url)(w, r)
	}).Methods("GET")
	router.HandleFunc("/update-fixture/{fixtureId}", func(w http.ResponseWriter, r *http.Request) {
		url := "http://localhost:8123/update-fixture/" + mux.Vars(r)["fixtureId"]
		forwardRequest(url)(w, r)
	}).Methods("PUT")
	router.HandleFunc("/delete-fixture/{fixtureId}", func(w http.ResponseWriter, r *http.Request) {
		url := "http://localhost:8124/delete-fixture/" + mux.Vars(r)["fixtureId"]
		forwardRequest(url)(w, r)
	}).Methods("DELETE")
	router.HandleFunc("/get-fixtures", forwardRequest("http://localhost:8123/get-fixtures")).Methods("GET")
}
