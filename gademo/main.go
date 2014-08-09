package main

import (
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/send", sendHandler)
	http.ListenAndServe(":8080", nil)
}

type Page struct {
	Title string
	Label string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Place an Order!", Label: "Demo"}
	t, _ := template.ParseFiles("template1.html")
	t.Execute(w, p)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {

}
