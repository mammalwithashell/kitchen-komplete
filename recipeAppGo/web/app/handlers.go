package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Pantry page handler
func pantryPage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Pantry Page")
}

// Recipe page handler
func createRecipePage(res http.ResponseWriter, req *http.Request) {
	// Sample Write to the database
	// Update this to use forms
	collection := client.Database("testdb").Collection("recipes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := Recipe{
		ID:    primitive.NewObjectID(),
		Name:  "Fish",
		Class: "Entree",
	}
	result, _ := collection.InsertOne(ctx, r)
	fmt.Fprint(res, result.InsertedID)
}

// display recipes saved in database
func readRecipePage(res http.ResponseWriter, req *http.Request) {

}

// update recipe saved in database
func updateRecipePage(res http.ResponseWriter, req *http.Request) {
	// Update in the CRUD

}

// delete recipe saved in database
func deleteRecipePage(res http.ResponseWriter, req *http.Request) {
	// Delete from the CRUD

}

//Struct to pass variable to support template
type supportStruct struct {
	ID      primitive.ObjectID `bson:"id,omitempty"`
	Email   string             `bson:"email,omitempty"`
	Subject string             `bson:"subject,omitempty"`
	Message string             `bson:"message,omitempty"`
}

// Support page handler
func supportPage(res http.ResponseWriter, req *http.Request) {
	// Let users report issues
	t, err1 := template.ParseFiles("./ui/html/support.page.tmpl")
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
	details := supportStruct{
		ID:      primitive.NewObjectID(),
		Email:   req.FormValue("email"),
		Subject: req.FormValue("subject"),
		Message: req.FormValue("message"),
	}

	collection := client.Database("testdb").Collection("support")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection.InsertOne(ctx, details)
	fmt.Fprint(res, "Thanks for your response!")
}
