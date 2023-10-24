package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongo(collectionName string) (*mongo.Collection, error) {
	err := godotenv.Load("../backend.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}
	mongoUser := os.Getenv("MONGO_USER")
	mongoPass := os.Getenv("MONGO_PASS")
	databaseName := os.Getenv("DATABASE_NAME")
	mongoUrl := "mongodb+srv://" + mongoUser + ":" + mongoPass + "@cms.drwtyxi.mongodb.net/" + databaseName + "?retryWrites=true&w=majority"

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoUrl).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database(mongoUser).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	collection := client.Database(databaseName).Collection(collectionName)
	return collection, nil
}
