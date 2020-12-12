package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Recipe struct for mongodb
type Recipe struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Category     string             `bson:"category,omitempty"`
	Ingredients  []string           `bson:"ingredients,omitempty"`
	PrepTime     string             `bson:"prepTime,omitempty"`
	Instructions []string           `bson:"instructions,omitempty"`
}

var t *template.Template
var pd Recipe

func servePage(res http.ResponseWriter, req *http.Request) {
	err := t.ExecuteTemplate(res, "index.page.gohtml", pd)
	if err != nil {
		fmt.Fprint(res, err)
	}
}

func main() {
	var err error
	t, err = template.ParseGlob("./ui/html/*")
	if err != nil {
		log.Println("Cannot parse templates:", err)
		os.Exit(-1)
	}
	pd = Recipe{
		Name:         "Noodles",
		Category:     "Entree",
		Ingredients:  []string{"Noodles", "otherstuff"},
		PrepTime:     "20 min",
		Instructions: []string{"cook", "it"},
	}
	http.HandleFunc("/", servePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
