package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/mypantry", pantryPage)
	http.HandleFunc("/support", supportPage)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
