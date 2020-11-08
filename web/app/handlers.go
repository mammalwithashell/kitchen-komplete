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
	"strings"
	"time"
	"unicode"

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
		ID:          primitive.NewObjectID(),
		Name:        req.FormValue("title"),
		Class:       req.FormValue("class"),
		Ingredients: strings.Split(req.FormValue("ingredients"), ","),
	}
	collection.InsertOne(ctx, r)
	fmt.Fprint(res, "Inserted a recipe into the database: ", r)
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
	t, err1 := app.templates.ParseFiles("./ui/html/support.page.gohtml")
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
	fmt.Fprint(res, err)

	// Load login template on anything other than a post request
	if req.Method != http.MethodPost {
		err := app.templates.ExecuteTemplate(res, "login.page.gohtml", nil)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	// Hash password before storing
	usr := models.User{
		UserID: req.FormValue("username"),
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
		fmt.Fprintln(res, "Username or Password aren't correct, please try again.")
		return
	}
	userPass := []byte(req.FormValue("passwd"))
	dbPass := []byte(temp1.HashedPassword)
	if passErr := bcrypt.CompareHashAndPassword(dbPass, userPass); passErr != nil {
		// If not the right password
		fmt.Fprintln(res, "Wrong Password")
		fmt.Fprintln(res, "DbPass: ", dbPass)
		fmt.Fprintln(res, "userPass: ", userPass)
		return
	}

	// Should be authenticated at this point on
	temp1.Authenticated = true
	session.Values["user"] = temp1
	err = session.Save(req, res)
	if err != nil {
		// Handle session save error
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/profile{}.html", http.StatusFound)
}

// Password function that ensures rigorous passwords
func Password(pass string) bool {
	var (
		upp, low, num, sym bool
		tot                uint8
	)

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			return false
		}
	}

	if !upp || !low || !num || !sym || tot < 8 {
		return false
	}

	return true
}

// Handler function for user registration
func (app application) registerHandler(res http.ResponseWriter, req *http.Request) {
	// Let users report issues
	t, err1 := app.templates.ParseFiles("./ui/html/register.page.gohtml")
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}

	// Give back login template on anything other than post
	if req.Method != http.MethodPost {
		err2 := t.Execute(res, nil)
		if err2 != nil {
			fmt.Println(err2.Error())
		}
		return
	}

	// Check if password meets all requirements
	if Password(req.FormValue("passwd")) == false {
		fmt.Fprint(res, "Password must contain all of the following: [A-Z], [0-9], [symbols or punctuations]")
		// Redirect to the same page again and prompt user.
		http.Redirect(res, req, "/login", http.StatusFound)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.FormValue("passwd")), bcrypt.MinCost)
	newusr := models.User{
		ID:             primitive.NewObjectID(),
		UserID:         req.FormValue("username"),
		Name:           req.FormValue("name"),
		Email:          req.FormValue("email"),
		HashedPassword: string(hash),
	}

	// Add co
	collection := app.client.Database("testdb").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var temp models.User
	err := collection.FindOne(ctx, bson.M{"user_id": newusr.UserID}).Decode(&temp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			collection.InsertOne(ctx, newusr)
			fmt.Fprint(res, "Inserted: ", newusr)
			return
		}
		// Redirect to profile page
	}

	fmt.Fprint(res, "Thanks for your response!")
}
