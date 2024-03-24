package config

import (
	"context"
	"log"

	//"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func init() {
	InitDB()
}

func InitDB() {
	// uri := os.Getenv("MONGODB_URI")
	// if uri == "" {
	// 	log.Fatal("MONGODB_URI environment variable is not set")
	// }
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	//Check connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	log.Println("Connected to MongoDB")

	DB = client.Database("golang_db") // Change the database name accordingly
}
