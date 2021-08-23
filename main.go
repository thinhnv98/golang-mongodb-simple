package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"time"
)

// Person represents a person document in MongoDB
type Person struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Age         int                `json:"age,omitempty" bson:"age,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
}

func homeLink(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Welcome home!")
}

var collection *mongo.Collection

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/admin")
	client, _ := mongo.Connect(ctx, clientOptions)

	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		print(err.Error())
	} else {
		fmt.Println("Ping succeed")
	}

	collection = client.Database("admin").Collection("people")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/people", createPerson).Methods("POST")
	fmt.Println("Running: 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createPerson(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, person)
	if err != nil {
		json.NewEncoder(response).Encode(err)
		return
	}
	id := result.InsertedID
	person.ID = id.(primitive.ObjectID)
	json.NewEncoder(response).Encode(person)
}
