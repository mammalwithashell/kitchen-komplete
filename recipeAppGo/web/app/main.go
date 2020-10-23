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

	/*//Prep flags to serve static pages
	var dir string
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()*/

	// Create multiplexer and handle routes
	router := mux.NewRouter()
	router.HandleFunc("/mypantry", pantryPage)
	router.HandleFunc("/myrecipes", recipePage)
	router.HandleFunc("/support", supportPage)

	// Serve static pages
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../static")))

	fmt.Println("Local Server running on port " + port)
	fmt.Println("http://localhost:8080")
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
