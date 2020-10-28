package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	//"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Declare mongo client globally
var client *mongo.Client

// Cookies/Session
// SSl

// Main function for serving http pages
func main() {
	port := ":8080"
	fmt.Println("Starting Kitchen Komplete application...")

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
	router.HandleFunc("/add_recipe", createRecipePage)
	router.HandleFunc("/mypantry", readRecipePage)
	router.HandleFunc("/update-recipe{number}/", updateRecipePage)
	router.HandleFunc("/support", supportPage)

	// Serve static pages
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/static")))

	// Start server for debug
	fmt.Println("Local Server running on port " + port)
	fmt.Println("http://localhost" + port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
