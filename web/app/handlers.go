package main

/*
Logic for the handler functions
*/
import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"mammal.shell/kitchenKomplete/models"
)

// struct for passsing data to the html templates
type templateData struct {
	Rec           []models.Recipe
	Authenticated bool
	Path          string
	User          string
	UserObj       models.User
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI, r.Method)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// Pantry page handler
func (app application) allRecipePage(res http.ResponseWriter, req *http.Request) {

	// Get userId from cookie or login

	// Find this users recipe collection
	collection := app.client.Database("recipes").Collection("test_user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	arr := []models.Recipe{}
	cur, _ := collection.Find(ctx, bson.D{})
	defer cur.Close(ctx)

	// Loop through found documents in collection
	for cur.Next(ctx) {
		result := models.Recipe{}
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		arr = append(arr, result)
	}

	// Display the items in arr
	fmt.Fprintf(res, "Pantry Page")
	for _, i := range arr {
		fmt.Fprintf(res, i.Name)
	}
}

func (app application) removeRecipe(res http.ResponseWriter, req *http.Request) {
	// Find this users recipe collection
	collection := app.client.Database("recipes").Collection("test_user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	arr := []models.Recipe{}
	cur, _ := collection.Find(ctx, bson.D{})
	defer cur.Close(ctx)

	// Loop through found documents in collection
	for cur.Next(ctx) {
		result := models.Recipe{}
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		arr = append(arr, result)
	}

	// Display the items in arr
	fmt.Fprintf(res, "Pantry Page")
	for _, i := range arr {
		fmt.Fprintf(res, i.Name)
	}
}

// Recipe page handler
func (app *application) createRecipe(res http.ResponseWriter, req *http.Request) {
	// Sample Write to the database
	// Find out who the user is
	// Pull the user info from database/cookie
	// Add
	// Update this to use forms
	var out templateData
	session, _ := app.store.Get(req, "cookie-name")
	sessionUser := session.Values["user"]
	if session.Values["authenticated"] != nil {
		out.Authenticated = session.Values["authenticated"].(bool)
	}

	if req.Method == http.MethodGet {
		app.templates.ExecuteTemplate(res, "add_recipe.page.gohtml", out)
		return
	}

	// If the recipe is set to public, a
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	usrRecipeCollection := app.client.Database("recipes").Collection(sessionUser.(string))

	r := models.Recipe{
		ID:           primitive.NewObjectID(),
		Name:         req.FormValue("name"),
		Category:     req.FormValue("category"),
		Ingredients:  strings.Split(req.FormValue("ingredients"), ","),
		PrepTime:     req.FormValue("preptime"),
		Instructions: strings.Split(req.FormValue("instructions"), ","),
		OwnerID:      sessionUser.(string),
		Public:       req.FormValue("privacy"),
	}
	usrRecipeCollection.InsertOne(ctx, r)
	app.templates.ExecuteTemplate(res, "add_recipe.page.gohtml", out)
}

// display recipes saved in database
func (app application) readRecipe(res http.ResponseWriter, req *http.Request) {
	// Maybe a table of all recipes given by user
}

// update recipe saved in database
func (app *application) recipe(res http.ResponseWriter, req *http.Request) {
	// Update in the CRUD
	var out templateData
	session, _ := app.store.Get(req, "cookie-name")
	sessionUser := session.Values["user"]
	if session.Values["authenticated"] != nil {
		out.Authenticated = session.Values["authenticated"].(bool)
	}

	if req.Method == http.MethodGet {
		app.templates.ExecuteTemplate(res, "recipe.page.gohtml", out)
		// fmt.Fprintf(res, "hello")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	usrRecipeCollection := app.client.Database("recipes").Collection(sessionUser.(string))

	r := models.Recipe{
		ID:           primitive.NewObjectID(),
		Name:         req.FormValue("name"),
		Category:     req.FormValue("category"),
		Ingredients:  strings.Split(req.FormValue("ingredients"), ","),
		PrepTime:     req.FormValue("preptime"),
		Instructions: strings.Split(req.FormValue("instructions"), ","),
		OwnerID:      sessionUser.(string),
		Public:       req.FormValue("privacy"),
	}
	usrRecipeCollection.InsertOne(ctx, r)
	app.templates.ExecuteTemplate(res, "profile_recipe.page.gohtml", out)
}

func (app *application) edit1(res http.ResponseWriter, req *http.Request) {
	var out templateData
	session, _ := app.store.Get(req, "cookie-name")
	out.User = session.Values["user"].(string)
	if session.Values["authenticated"] != nil {
		out.Authenticated = session.Values["authenticated"].(bool)
	}
	err := app.templates.ExecuteTemplate(res, "edit.page.gohtml", out)
	if err != nil {
		fmt.Fprint(res, err)
	}
}

func (app *application) editRecipe(res http.ResponseWriter, req *http.Request) {
	//setup
	var out templateData
	session, _ := app.store.Get(req, "cookie-name")
	out.User = session.Values["user"].(string)

	if out.Authenticated = auth(session, nil); !out.Authenticated {
		http.Redirect(res, req, "/", http.StatusForbidden)
		return
	}

	var recipeForUpdate models.Recipe

	// Connect to database and collection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	usrRecipeCollection := app.client.Database("recipes").Collection(out.User)
	u, _ := url.Parse(req.URL.String())
	id := strings.Split(u.Path, "/")[2]
	objID, _ := primitive.ObjectIDFromHex(id)
	fmt.Println("User:", out.User)
	err := usrRecipeCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&recipeForUpdate)
	fmt.Println(recipeForUpdate)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}

	// fill the templateData struct
	out.Rec = append(out.Rec, recipeForUpdate)

	if req.Method != http.MethodPost {
		err = app.templates.ExecuteTemplate(res, "edit.page.gohtml", out)
		if err != nil {
			fmt.Fprint(res, err)
		}
		return
	}

	// Fill recipe struct from user input
	r := models.Recipe{
		ID:           objID,
		Name:         req.FormValue("name"),
		Category:     req.FormValue("category"),
		Ingredients:  strings.Split(req.FormValue("ingredients"), ","),
		PrepTime:     req.FormValue("preptime"),
		Instructions: strings.Split(req.FormValue("instructions"), ","),
		OwnerID:      out.User,
		Public:       req.FormValue("privacy"),
	}

	//implemet
	filter := bson.D{{"_id", r.ID}}
	// update := bson.D{{"$set", bson.D{{"email", "newemail@example.com"}}}}

	result, err := usrRecipeCollection.UpdateOne(ctx, filter, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return
	}
}

// delete recipe saved in database
func (app application) deleteRecipe(res http.ResponseWriter, req *http.Request) {
	// Delete from the CRUD
	//vars := mux.Vars(req)
}

// Function to display a user's pantry
func (app application) pantry(res http.ResponseWriter, req *http.Request) {

}

func (app application) updatePantry(res http.ResponseWriter, req *http.Request) {

}

// Support page handler
func (app application) support(res http.ResponseWriter, req *http.Request) {
	// Let users report issues
	t, err1 := template.ParseFiles("./ui/html/support.page.gohtml")
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}

	// Give back support template on anything other than post
	if req.Method != http.MethodPost {
		err2 := t.Execute(res, nil)
		if err2 != nil {
			fmt.Println(err2.Error())
		}
		return
	}

	// Once the form is submitted, the route is hit with a post request
	// Create the support struct for database
	details := models.Support{
		ID:      primitive.NewObjectID(),
		Email:   req.FormValue("email"),
		Subject: req.FormValue("subject"),
		Message: req.FormValue("message"),
	}

	collection := app.client.Database("testdb").Collection("support")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection.InsertOne(ctx, details)
	fmt.Fprint(res, "Thanks for your response!")
	return
}

// Handler funcion for user login
func (app *application) login(res http.ResponseWriter, req *http.Request) {
	var out templateData
	session, err := app.store.Get(req, "cookie-name")

	// Load login template on anything other than a post request
	t, _ := template.ParseFiles("./ui/html/login.page.gohtml")
	if req.Method != http.MethodPost {
		app.templates.ExecuteTemplate(res, "login.page.gohtml", out)
		return
	}

	// Hash password before storing
	usr := models.User{
		UserID: req.FormValue("username"),
		Errors: make(map[string]string),
	}

	// Select the collection to be queried
	collection := app.client.Database("Users").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Look for a match in the database
	var newUser models.User
	err = collection.FindOne(ctx, bson.M{"user_id": usr.UserID}).Decode(&newUser)
	if err != nil {
		// Handle error if user_id is not in database
		res.WriteHeader(http.StatusInternalServerError)
		usr.Errors["Login"] = "Username or Password isn't correct."
		app.templates.ExecuteTemplate(res, "/login.page.gohtml", out)
		return
	}
	userPass := []byte(req.FormValue("passwd"))
	dbPass := []byte(newUser.HashedPassword)
	if passErr := bcrypt.CompareHashAndPassword(dbPass, userPass); passErr != nil {
		usr.Errors["Login"] = "Username or Password isn't correct."
		t.Execute(res, usr)
		return
	}

	// Should be authenticated at this point on
	// newUser.Authenticated = true
	session.Values["authenticated"] = true
	session.Values["user"] = newUser.UserID
	// Save session cookie
	err = session.Save(req, res)
	if err != nil {
		// Handle session save error
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/profile", http.StatusFound)
}

// Logic for logout
func (app application) logout(res http.ResponseWriter, req *http.Request) {
	session, _ := app.store.Get(req, "cookie-name")
	session.Values["authenticated"] = false
	session.Values["user"] = nil
	session.Save(req, res)
	http.Redirect(res, req, "/", http.StatusFound)
}

// Handler function for user registration
func (app *application) register(res http.ResponseWriter, req *http.Request) {
	var out templateData
	session, _ := app.store.Get(req, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		out.Authenticated = true
	}

	// Give back register template on anything other than post
	if req.Method != http.MethodPost {
		err2 := app.templates.ExecuteTemplate(res, "register.page.gohtml", out)
		if err2 != nil {
			fmt.Println(err2.Error())
		}
		return
	}
	out.UserObj = models.User{
		ID:     primitive.NewObjectID(),
		UserID: req.FormValue("username"),
		Name:   req.FormValue("name"),
		Email:  req.FormValue("email"),
		Errors: make(map[string]string),
	}

	// Check if Name and Username were valid
	if req.FormValue("name") == "" {
		out.UserObj.Errors["Name"] = "This field cannot be empty"
		app.templates.ExecuteTemplate(res, "register.page.gohtml", out)
		return
	}
	// Check if username field is not empty
	if req.FormValue("username") == "" {
		out.UserObj.Errors["User"] = "Please enter a Username"
		app.templates.ExecuteTemplate(res, "register.page.gohtml", out)
		return
	}
	// Check if email & password was valid
	if out.UserObj.Validate() == false || out.UserObj.Password(req.FormValue("passwd")) == false {
		app.templates.ExecuteTemplate(res, "register.page.gohtml", out)
		return
	}
	// Check if Passwords match
	if req.FormValue("cfm_passwd") != req.FormValue("passwd") {
		out.UserObj.Errors["Password"] = "Please make your passwords match"
		app.templates.ExecuteTemplate(res, "register.page.gohtml", out)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.FormValue("passwd")), bcrypt.MinCost)
	out.UserObj.HashedPassword = string(hash)

	// Access database
	collection := app.client.Database("Users").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if username is already in database
	var temp models.User
	err := collection.FindOne(ctx, bson.M{"user_id": out.UserObj.UserID}).Decode(&temp)
	if err != nil {
		// if username doesn't exist
		if err == mongo.ErrNoDocuments {
			collection.InsertOne(ctx, out.UserObj)
			// Redirect to profile page
			session.Values["user"] = out.UserObj.UserID
			session.Values["authenticated"] = true
			session.Save(req, res)
			http.Redirect(res, req, "/profile", http.StatusFound)
			return
		}
	} else {

		out.UserObj.Errors["User"] = "Username Already exists!"
		app.templates.ExecuteTemplate(res, "register.page.gohtml", out)
		return
	}
}

// Show recipes
func (app *application) showRecipe(res http.ResponseWriter, req *http.Request) {
	var out templateData
	out.Authenticated = auth(app.store.Get(req, "cookie-name"))

	db := app.client.Database("recipes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, _ := db.ListCollectionNames(ctx, bson.M{}, nil)

	for _, name := range cursor {
		col := db.Collection(name)
		cur, _ := col.Find(ctx, bson.M{"public": "Public"})

		// Loop through found documents in collection
		for cur.Next(ctx) {
			result := models.Recipe{}
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			out.Rec = append(out.Rec, result)
		}
	}
	app.templates.ExecuteTemplate(res, "index.page.gohtml", out)
}

// Show profile
func (app *application) profile(res http.ResponseWriter, req *http.Request) {
	var out templateData
	session, _ := app.store.Get(req, "cookie-name")
	if out.Authenticated = auth(session, nil); !out.Authenticated {
		http.Redirect(res, req, "/", http.StatusForbidden)
		return
	}

	// Isolate profile in url
	u, _ := url.Parse(req.URL.String())
	out.Path = u.Path
	out.User = session.Values["user"].(string)

	userRecipeCollection := app.client.Database("recipes").Collection(out.User)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, _ := userRecipeCollection.Find(ctx, bson.D{})
	defer cur.Close(ctx)

	// Store all recipes in out.Rec
	if err := cur.All(ctx, &out.Rec); err != nil {
		log.Fatal(err)
	}
	app.templates.ExecuteTemplate(res, "profile.page.gohtml", out)
}

// Show about page
func (app *application) about(res http.ResponseWriter, req *http.Request) {
	var out templateData
	out.Authenticated = auth(app.store.Get(req, "cookie-name"))
	u, _ := url.Parse(req.URL.String())
	out.Path = u.Path
	app.templates.ExecuteTemplate(res, "about.page.gohtml", out)
}
