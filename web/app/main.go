package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"mammal.shell/kitchenKomplete/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	port := ":8080"
	fmt.Println("Starting Kitchen Komplete application...")
	var app application

	// Connect to mongodb client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	app.client, err = mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://james:VBQHfudvGQyoSyyf@cluster0.ehf5d.mongodb.net/sample_supplies?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	defer app.client.Disconnect(ctx)

	// Create multiplexer and handle routes
	app.router = mux.NewRouter()
	app.templates = template.Must(template.ParseGlob("./ui/html/*.gohtml"))

	app.router.HandleFunc("/all_recipes", app.allRecipePage)
	app.router.HandleFunc("/add_recipe", app.createRecipePage)
	app.router.HandleFunc("/myrecipes", app.readRecipePage)
	app.router.HandleFunc("/update-recipe{number}/", app.updateRecipePage)
	app.router.HandleFunc("/remove-recipe{_id}", app.deleteRecipeHandlerFunc)
	app.router.HandleFunc("/support", app.supportPage)
	app.router.HandleFunc("/mypantry", app.pantryHandler)
	app.router.HandleFunc("/login", app.loginHandler)
	app.router.HandleFunc("/register", app.registerHandler)
	// serve static pages
	app.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/static")))

	// Start server for debug
	fmt.Println("Local Server running on port " + port)
	fmt.Println("http://localhost" + port)

	if err := http.ListenAndServe(port, app.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// A struct to represent key things in the application
type application struct {
	router    *mux.Router
	client    *mongo.Client         // Declare mongo client globally
	store     *sessions.CookieStore // Cookies/Session
	templates *template.Template
	// SSl

}

// Function to get user cookie.
func getUser(s *sessions.Session) models.User {
	val := s.Values["user"]
	var user = models.User{}
	user, ok := val.(models.User)
	if !ok {
		return models.User{Authenticated: false}
	}
	return user
}
