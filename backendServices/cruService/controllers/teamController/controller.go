package teamController

import (
	"CMS/backendServices/config"
	models "CMS/backendServices/models"
	responses "CMS/backendServices/response"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var teamCollection, _ = config.ConnectToMongo("Teams")
var validate = validator.New()

func CreateTeamHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var team models.Team
		defer cancel()
		if err := json.NewDecoder(request.Body).Decode(&team); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.TeamResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		if validationErr := validate.Struct(&team); validationErr != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.TeamResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		newPlayer := models.Team{
			Id:            primitive.NewObjectID(),
			Name:          team.Name,
			Members:       team.Members,
			Captain:       team.Captain,
			PlayingXI:     team.PlayingXI,
			MatchesPlayed: team.MatchesPlayed,
			MatchesWon:    team.MatchesWon,
			MatchesLost:   team.MatchesLost,
			MatchesTied:   team.MatchesTied,
			MatchPoints:   team.MatchPoints,
		}

		result, err := teamCollection.InsertOne(ctx, newPlayer)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.TeamResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		response.WriteHeader(http.StatusCreated)
		jsonResponse := responses.TeamResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(response).Encode(jsonResponse)
	}
}

func GetTeamHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		params := mux.Vars(request)
		teamId := params["teamId"]
		var team models.Team
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(teamId)

		err := teamCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&team)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.TeamResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		response.WriteHeader(http.StatusOK)
		jsonResponse := responses.TeamResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": team}}
		json.NewEncoder(response).Encode(jsonResponse)
	}
}

func UpdateTeamHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		params := mux.Vars(request)
		teamId := params["teamId"]
		var team models.Team
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(teamId)

		//validate the request body
		if err := json.NewDecoder(request.Body).Decode(&team); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.TeamResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&team); validationErr != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.TeamResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		update := bson.M{
			"matchPoints":   team.MatchPoints,
			"matchesPlayed": team.MatchesPlayed,
			"matchesWon":    team.MatchesWon,
			"matchesLost":   team.MatchesLost,
			"matchesTied":   team.MatchesTied,
		}

		result, err := teamCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		//get updated player details
		var updatedTeam models.Team
		if result.MatchedCount == 1 {
			err := teamCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedTeam)

			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				jsonResponse := responses.TeamResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(response).Encode(jsonResponse)
				return
			}
		}

		response.WriteHeader(http.StatusOK)
		jsonResponse := responses.TeamResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedTeam}}
		json.NewEncoder(response).Encode(jsonResponse)
	}
}

func GetAllTeamsHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var teams []models.Team
		defer cancel()

		results, err := teamCollection.Find(ctx, bson.M{})

		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.TeamResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleTeam models.Team
			if err = results.Decode(&singleTeam); err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				jsonResponse := responses.TeamResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(response).Encode(jsonResponse)
			}

			teams = append(teams, singleTeam)
		}

		response.WriteHeader(http.StatusOK)
		jsonResponse := responses.TeamResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"teams": teams}}
		json.NewEncoder(response).Encode(jsonResponse)
	}
}
