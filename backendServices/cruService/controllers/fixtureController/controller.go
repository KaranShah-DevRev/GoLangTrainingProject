package fixturecontroller

import (
	"CMS/backendServices/config"
	"CMS/backendServices/models"
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

var fixtureCollection, err = config.ConnectToMongo("Fixtures")
var validate = validator.New()

func CreateFixtureHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.FixtureResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var fixture models.Fixture
		defer cancel()
		if err := json.NewDecoder(request.Body).Decode(&fixture); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.FixtureResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		if validationErr := validate.Struct(&fixture); validationErr != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.FixtureResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		newFixture := models.Fixture{
			Id:           primitive.NewObjectID(),
			MatchNumber:  fixture.MatchNumber,
			TeamA:        fixture.TeamA,
			TeamB:        fixture.TeamB,
			MatchDate:    fixture.MatchDate,
			Venue:        fixture.Venue,
			IsFinalMatch: fixture.IsFinalMatch,
		}

		result, err := fixtureCollection.InsertOne(ctx, newFixture)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.FixtureResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		response.WriteHeader(http.StatusCreated)
		jsonResponse := responses.FixtureResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(response).Encode(jsonResponse)

	}
}

func GetFixtureHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.FixtureResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		params := mux.Vars(request)
		fixtureId := params["fixtureId"]
		var fixture models.Fixture
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(fixtureId)

		err := fixtureCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&fixture)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.FixtureResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		response.WriteHeader(http.StatusOK)
		jsonResponse := responses.FixtureResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": fixture}}
		json.NewEncoder(response).Encode(jsonResponse)
	}
}

func GetAllFixturesHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.FixtureResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var fixtures []models.Fixture
		defer cancel()

		results, err := fixtureCollection.Find(ctx, bson.M{})

		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.FixtureResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleFixture models.Fixture
			if err = results.Decode(&singleFixture); err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				jsonResponse := responses.FixtureResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(response).Encode(jsonResponse)
			}

			fixtures = append(fixtures, singleFixture)
		}

		response.WriteHeader(http.StatusOK)
		jsonResponse := responses.FixtureResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"fixtures": fixtures}}
		json.NewEncoder(response).Encode(jsonResponse)
	}
}

func UpdateFixtureHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.FixtureResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		params := mux.Vars(request)
		fixtureId := params["fixtureId"]
		var fixture models.Fixture
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(fixtureId)

		//validate the request body
		if err := json.NewDecoder(request.Body).Decode(&fixture); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.FixtureResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&fixture); validationErr != nil {
			response.WriteHeader(http.StatusBadRequest)
			jsonResponse := responses.FixtureResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		update := bson.M{
			"matchDate": fixture.MatchDate,
			"venue":     fixture.Venue,
		}

		result, err := fixtureCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			jsonResponse := responses.FixtureResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(response).Encode(jsonResponse)
			return
		}

		//get updated fixture details
		var updatedFixture models.Fixture
		if result.MatchedCount == 1 {
			err := fixtureCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedFixture)

			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				jsonResponse := responses.FixtureResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(response).Encode(jsonResponse)
				return
			}
		}

		response.WriteHeader(http.StatusOK)
		jsonResponse := responses.FixtureResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedFixture}}
		json.NewEncoder(response).Encode(jsonResponse)
	}
}
