package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func homePage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello")
}

func pantryPage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello")
}

type supportStruct struct {
	title string
}

func supportPage(res http.ResponseWriter, req *http.Request) {
	if t, err1 := template.ParseFiles("index.html"); err1 != nil {
		panic(err1.Error)
	} else {
		err2 := t.Execute(res, supportStruct{"Hello"})
		if err2 != nil {
			panic(err2.Error)
		}
	}

	fmt.Fprintf(res, "Hello")
}
