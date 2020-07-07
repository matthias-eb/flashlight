package controller

import (
	"fmt"
	"net/http"

	mw "github.com/matthias-eb/flashlight/app/middleware"
	db "github.com/matthias-eb/flashlight/app/model"
	st "github.com/matthias-eb/flashlight/app/structs"
)

//Preview responds to a Get Request to Root. It will then show the index Page with the newest Posts of all Users
func Preview(w http.ResponseWriter, r *http.Request) {
	var isAuthenticated bool
	mw.SetupSession(w, r)
	username, err := mw.CheckAuthentication(w, r)
	isAuthenticated = err == nil // isAuthenticated is true if the error was nil.
	data := st.Data{
		Title: "Flashlight",
		Error: nil,
	}
	images, err := db.GetAllImages(username)
	if err != nil {
		fmt.Printf("Error executing Templates: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isAuthenticated {
		data.User = ""
		fmt.Printf("User is not logged in.\n")
	} else {
		fmt.Printf("User %+v is logged in.\n", username)
		data.User = username
	}
	data.Images = images
	err = mw.Templ.ExecuteTemplate(w, "index.tmpl", data)
	if err != nil {
		fmt.Printf("Error executing Templates: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
