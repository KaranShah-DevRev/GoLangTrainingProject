package playerController

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

var playerCollection, err = config.ConnectToMongo("Players")
var validate = validator.New()

func CreatePlayerHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var player models.Player
		defer cancel()
		if err := json.NewDecoder(request.Body).Decode(&player); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.PlayerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		if validationErr := validate.Struct(&player); validationErr != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.PlayerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		newPlayer := models.Player{
			Id:            primitive.NewObjectID(),
			Name:          player.Name,
			Age:           player.Age,
			Role:          player.Role,
			BattingAvg:    player.BattingAvg,
			StrikeRate:    player.StrikeRate,
			Economy:       player.Economy,
			MatchesPlayed: player.MatchesPlayed,
			TotalRuns:     player.TotalRuns,
			TotalWickets:  player.TotalWickets,
		}

		result, err := playerCollection.InsertOne(ctx, newPlayer)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		response.WriteHeader(http.StatusCreated)
		jsonResponse := responses.PlayerResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(response).Encode(jsonResponse)

	}
}

func GetPlayerHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		params := mux.Vars(request)
		playerId := params["playerId"]
		var player models.Player
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(playerId)

		err := playerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&player)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		response.WriteHeader(http.StatusOK)
		jsonResponse := responses.PlayerResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": player}}
		json.NewEncoder(response).Encode(jsonResponse)
	}
}

func UpdatePlayerHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		params := mux.Vars(request)
		playerId := params["playerId"]
		var player models.Player
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(playerId)

		//validate the request body
		if err := json.NewDecoder(request.Body).Decode(&player); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.PlayerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&player); validationErr != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.PlayerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		update := bson.M{
			"Age":           player.Age,
			"BattingAvg":    player.BattingAvg,
			"StrikeRate":    player.StrikeRate,
			"Economy":       player.Economy,
			"MatchesPlayed": player.MatchesPlayed,
			"TotalRuns":     player.TotalRuns,
			"TotalWickets":  player.TotalWickets,
		}

		result, err := playerCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		//get updated player details
		var updatedPlayer models.Player
		if result.MatchedCount == 1 {
			err := playerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedPlayer)

			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				jsonResponse := responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(response).Encode(jsonResponse)
				return
			}
		}

		response.WriteHeader(http.StatusOK)
		jsonResponse := responses.PlayerResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedPlayer}}
		json.NewEncoder(response).Encode(jsonResponse)
	}
}

func GetAllPlayers() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var players []models.Player
		defer cancel()

		results, err := playerCollection.Find(ctx, bson.M{})

		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singlePlayer models.Player
			if err = results.Decode(&singlePlayer); err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				jsonResponse := responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(response).Encode(jsonResponse)
			}

			players = append(players, singlePlayer)
		}

		response.WriteHeader(http.StatusOK)
		jsonResponse := responses.PlayerResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"players": players}}
		json.NewEncoder(response).Encode(jsonResponse)
	}
}
