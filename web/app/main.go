package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mammal.shell/kitchenKomplete/models"
)

// Main function for serving http pages
var app application

// Init function runs before the main go code
func init() {

	// Secure Cookie stuff
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)
	app.store = sessions.NewCookieStore(authKeyOne, encryptionKeyOne)

	// Options for all the session cookies
	app.store.Options = &sessions.Options{
		MaxAge:   60 * 15, // 15 minutes in units of seconds
		HttpOnly: true,
	}

	// Register the user type with gob/encoding so it can be written as a session value
	gob.Register(models.User{})
}

func main() {
	//Uncomment to load environment variables from .env file
	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("Error loading .env file")
	}

	// Set environment variables for static ip service
	os.Setenv("HTTP_PROXY", os.Getenv("IPB_HTTP"))
	os.Setenv("HTTPS_PROXY", os.Getenv("IPB_HTTPS"))

	port := os.Getenv("PORT")
	dburi := os.Getenv("MONGODB_URI")
	fmt.Println("Starting Kitchen Komplete application...")

	// Connect to mongodb client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	app.client, err = mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	defer app.client.Disconnect(ctx)

	// Create multiplexer and handle routes
	app.router = mux.NewRouter()
	//app.templates = packr.NewBox("./ui/templates")
	//app.static = packr.NewBox("./ui/static")
	//app.templates = template.Must(template.ParseGlob("./ui/html/*.gohtml"))
	app.Routes()

	// serve static pages
	app.templates = template.Must(template.ParseGlob("./ui/html/*"))
	app.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/static/")))

	// Start server for debug
	fmt.Println("Local Server running on port " + port)
	fmt.Println("http://localhost:" + port)

	if err := http.ListenAndServe(":"+port, app.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// A struct to represent key things in the application
type application struct {
	router    *mux.Router           // Router
	client    *mongo.Client         // Declare mongo client globally
	store     *sessions.CookieStore // Cookies/Session
	templates *template.Template    // templates
	static    packr.Box
	// SSl

}

// Function to get userID from session cookie.
func getUser(s *sessions.Session) models.User {
	val := s.Values["user"]
	var user = models.User{}
	user, ok := val.(models.User)
	if !ok {
		return models.User{Authenticated: false}
	}
	return user
}
