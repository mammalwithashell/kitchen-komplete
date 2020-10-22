package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Declare mongo client globally
var client *mongo.Client

// Main function for serving http pages
func main() {
	port := ":8080"
	fmt.Println("Starting Kitchen Komplete application...")
	fmt.Println("Local Server running on port" + port)

	// Connect to mongodb client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://james:VBQHfudvGQyoSyyf@cluster0.ehf5d.mongodb.net/sample_supplies?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Create multiplexer and handle routes
	router := mux.NewRouter()
	router.HandleFunc("/", homePage)
	router.HandleFunc("/mypantry", pantryPage)
	router.HandleFunc("/myrecipes", recipePage)
	router.HandleFunc("/support", supportPage)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
