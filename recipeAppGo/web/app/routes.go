package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Function to adapt file path for templates
func templatePath(s string) string {
	return "C:\\Users\\james\\OneDrive\\School\\Fall 2020\\CSCE 3444 Software Engineering\\Project1\\CSCE-3444-Team-AIJNW\\recipeAppGo\\web\\template\\" + s
}

// Home page handler
func homePage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello")
}

// Pantry page handler
func pantryPage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Pantry Page")
}

// Recipe page handler
func recipePage(res http.ResponseWriter, req *http.Request) {
	collection := client.Database("testdb").Collection("recipes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, bson.D{
		{Key: "dishName", Value: "Lemon Pepper Wings"},
		{Key: "class", Value: "Entree"},
	})
	fmt.Fprint(res, result.InsertedID)
}

//Struct to pass variable to support template
type supportStruct struct {
	Title string
}

// Support page handler
func supportPage(res http.ResponseWriter, req *http.Request) {
	if t, err1 := template.ParseFiles(templatePath("index.html")); err1 != nil {
		fmt.Println(err1.Error())
	} else {
		err2 := t.Execute(res, supportStruct{"Yeet"})
		if err2 != nil {
			fmt.Println(err2.Error())
		}
	}
}
