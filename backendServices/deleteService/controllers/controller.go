package controllers

import (
	"CMS/backendServices/config"
	responses "CMS/backendServices/response"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var playerCollection, _ = config.ConnectToMongo("Players")
var teamCollection, _ = config.ConnectToMongo("Teams")
var fixtureCollection, _ = config.ConnectToMongo("Fixtures")

func DeletePlayer() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		playerId := params["playerId"]
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(playerId)

		result, err := playerCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount < 1 {
			rw.WriteHeader(http.StatusNotFound)
			response := responses.PlayerResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Player with specified ID not found!"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.PlayerResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Player successfully deleted!"}}
		json.NewEncoder(rw).Encode(response)
	}
}

func DeleteTeam() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		teamId := params["teamId"]
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(teamId)

		result, err := teamCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.TeamResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount < 1 {
			rw.WriteHeader(http.StatusNotFound)
			response := responses.TeamResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Team with specified ID not found!"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.TeamResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Team successfully deleted!"}}
		json.NewEncoder(rw).Encode(response)
	}
}

func DeleteFixture() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		teamId := params["fixtureId"]
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(teamId)

		result, err := fixtureCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.FixtureResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount < 1 {
			rw.WriteHeader(http.StatusNotFound)
			response := responses.FixtureResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Fixture with specified ID not found!"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.FixtureResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Fixture successfully deleted!"}}
		json.NewEncoder(rw).Encode(response)
	}
}
