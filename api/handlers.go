package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"myproject/config"
	"myproject/models"
	"myproject/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbs        *mongo.Database
	collection *mongo.Collection
	validate   = validator.New()
)

func init() {
	dbs = config.DB
	collection = dbs.Collection("users")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(user); err != nil {
		// Validation failed
		var validationErrors []utils.ErrorResponse
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, utils.ErrorResponse{
				Field: err.Field(),
				Error: err.Tag(),
			})
		}

		// Encode the validation errors as JSON and write to the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErrors)
		return

	}
	log.Printf("Req: %+v", user)
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Result: %v", result)

	json.NewEncoder(w).Encode(result)

}

func GetOneUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the request URL parameters
	vars := mux.Vars(r)
	userID := vars["id"]

	// Check if the user ID is provided
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// Define a filter to find the user by ID
	filter := bson.D{{Key: "_id", Value: objID}}
	log.Println(filter)
	// Perform the find operation to get the user
	var user models.User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		// If user not found or any other error occurs, return an error response
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Error querying user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Set the response header and encode the user as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func FindAll(w http.ResponseWriter, r *http.Request) {
	filter := bson.D{{}}
	// Perform the find operation
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Error querying users: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor and decode documents
	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Printf("Error decoding user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		log.Printf("Error iterating cursor: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response header and encode the users as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
