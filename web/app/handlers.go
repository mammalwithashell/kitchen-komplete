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
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"mammal.shell/kitchenKomplete/models"
)

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

// Recipe page handler
func (app application) createRecipePage(res http.ResponseWriter, req *http.Request) {
	// Sample Write to the database
	// Update this to use forms
	t, err1 := template.ParseFiles("./ui/html/addrecipe.page.gohtml")
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

	collection := app.client.Database("recipes").Collection("test_user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r := models.Recipe{
		ID:           primitive.NewObjectID(),
		Name:         req.FormValue("title"),
		Category:     req.FormValue("category"),
		Ingredients:  strings.Split(req.FormValue("ingredients"), ","),
		PrepTime:     req.FormValue("preptime"),
		Instructions: strings.Split(req.FormValue("instructions"), ","),
	}
	collection.InsertOne(ctx, r)
	http.Redirect(res, req, "/add_recipe.html", http.StatusFound)
}

// display recipes saved in database
func (app application) readRecipePage(res http.ResponseWriter, req *http.Request) {
	// Maybe a table of all recipes given by user
}

// update recipe saved in database
func (app application) updateRecipePage(res http.ResponseWriter, req *http.Request) {
	// Update in the CRUD

}

// delete recipe saved in database
func (app application) deleteRecipeHandlerFunc(res http.ResponseWriter, req *http.Request) {
	// Delete from the CRUD
	//vars := mux.Vars(req)
}

// Function to display a user's pantry
func (app application) pantryHandler(res http.ResponseWriter, req *http.Request) {

}

// Support page handler
func (app application) supportPage(res http.ResponseWriter, req *http.Request) {
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

// Hash Passwords
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Check passwords
func checkHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Check Hash Err:", err)
		return false
	}
	return true
}

// Handler funcion for user login
func (app application) loginHandler(res http.ResponseWriter, req *http.Request) {
	session, err := app.store.Get(req, "cookie-name")
	// Load login template on anything other than a post request
	t, _ := template.ParseFiles("./ui/html/login.page.gohtml")
	if req.Method != http.MethodPost {
		t.Execute(res, nil)
		return
	}

	// Hash password before storing
	usr := models.User{
		UserID: req.FormValue("username"),
		Errors: make(map[string]string),
	}

	// Select the collection to be queried
	collection := app.client.Database("testdb").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Look for a match in the database
	var temp1 models.User
	err = collection.FindOne(ctx, bson.M{"user_id": usr.UserID}).Decode(&temp1)
	if err != nil {
		// Handle error if user_id is not in database
		res.WriteHeader(http.StatusInternalServerError)
		usr.Errors["Login"] = "Username or Password isn't correct."
		t.Execute(res, usr)
		return
	}
	userPass := []byte(req.FormValue("passwd"))
	dbPass := []byte(temp1.HashedPassword)
	if passErr := bcrypt.CompareHashAndPassword(dbPass, userPass); passErr != nil {
		usr.Errors["Login"] = "Username or Password isn't correct."
		t.Execute(res, usr)
		return
	}

	// Should be authenticated at this point on
	// temp1.Authenticated = true
	session.Values["authenticated"] = true
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
func (app application) logoutHandler(res http.ResponseWriter, req *http.Request) {
	session, _ := app.store.Get(req, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(req, res)
	http.Redirect(res, req, "/", http.StatusFound)
}

// Handler function for user registration
func (app application) registerHandler(res http.ResponseWriter, req *http.Request) {
	t, err1 := template.ParseFiles("./ui/html/register.page.gohtml")
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}

	// Give back register template on anything other than post
	if req.Method != http.MethodPost {
		err2 := t.Execute(res, nil)
		if err2 != nil {
			fmt.Println(err2.Error())
		}
		return
	}
	newusr := models.User{
		ID:     primitive.NewObjectID(),
		UserID: req.FormValue("username"),
		Name:   req.FormValue("name"),
		Email:  req.FormValue("email"),
		Errors: make(map[string]string),
	}

	// Check if Name and Username were valid
	if req.FormValue("name") == "" {
		newusr.Errors["Name"] = "Please enter a name"
		t.Execute(res, newusr)
		return
	}
	// Check if username field is not empty
	if req.FormValue("username") == "" {
		newusr.Errors["User"] = "Please enter a Username"
		t.Execute(res, newusr)
		return
	}
	// Check if email was valid
	if newusr.Validate() == false || newusr.Password(req.FormValue("passwd")) == false {
		t.Execute(res, newusr)
		return
	}
	// Check if Passwords match
	if req.FormValue("cfm_passwd") != req.FormValue("passwd") {
		newusr.Errors["Password"] = "Please make your passwords match"
		t.Execute(res, newusr)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.FormValue("passwd")), bcrypt.MinCost)
	newusr.HashedPassword = string(hash)

	// Access database
	collection := app.client.Database("testdb").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if username is already in database
	var temp models.User
	err := collection.FindOne(ctx, bson.M{"user_id": newusr.UserID}).Decode(&temp)
	if err != nil {
		// if username doesn't exist
		if err == mongo.ErrNoDocuments {
			collection.InsertOne(ctx, newusr)
			// Redirect to profile page
			http.Redirect(res, req, "/profile.html", http.StatusFound)
			return
		}
	} else {
		newusr.Errors["User"] = "Username Already exists!"
		t.Execute(res, newusr)
		return
	}
}

type showRecipes struct {
	Rec           []models.Recipe
	Authenticated bool
}

// Show recipes
func (app *application) showHandler(res http.ResponseWriter, req *http.Request) {
	var out showRecipes
	session, _ := app.store.Get(req, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		out.Authenticated = true
	}

	db := app.client.Database("recipes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, _ := db.ListCollectionNames(ctx, bson.M{}, nil)

	for _, name := range cursor {
		col := db.Collection(name)
		cur, _ := col.Find(ctx, bson.D{})

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
	file, _ := os.Open("error.txt")
	if err := app.templates.ExecuteTemplate(res, "index.page.gohtml", out); err != nil {
		fmt.Fprint(file, "Error: ", err)
	}
	fmt.Fprint(file, "OK")
}

// Show profile
func (app *application) profileHandler(res http.ResponseWriter, req *http.Request) {
	var out showRecipes
	session, _ := app.store.Get(req, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		out.Authenticated = true
	}
	app.templates.ExecuteTemplate(res, "profile.page.gohtml", nil)
}
