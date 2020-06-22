package controller

import (
	"net/http"
)

//Login is for logging in a User. It can receive a POST or a GET Method.
func Login(w http.ResponseWriter, r *http.Request) {

	type Data struct {
		Title string
		Error []string
	}

	if r.Method == "POST" {
		http.Error(w, "Not implemented yet", http.StatusInternalServerError)
	} else if r.Method == "GET" {
		tmpl.ExecuteTemplate(w, "login.tmpl", Data{Title: "Flashlight Login", Error: nil})
	}
}
