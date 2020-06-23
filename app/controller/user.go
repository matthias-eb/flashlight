package controller

import (
	"fmt"
	"net/http"

	mw "github.com/matthias-eb/flashlight/app/middleware"
	db "github.com/matthias-eb/flashlight/app/model"
)

//Login is for logging in a User. It can receive a POST or a GET Method.
func Login(w http.ResponseWriter, r *http.Request) {
	mw.SetupSession(w, r)
	defer mw.SaveSession(w, r)

	fmt.Println("Request for /login coming in: " + r.Method)

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := db.GetUser(username)

		if err == nil {
			err = mw.AuthenticateUser(username, password, user.Password)
		}

		if err != nil {
			mw.Templ.ExecuteTemplate(w, "login.tmpl", Data{Title: "Flashlight Login", Error: []string{"Benutzername oder Password waren falsch"}})
			return
		}

		http.Redirect(w, r, "/", http.StatusOK)

	} else if r.Method == "GET" {
		mw.Templ.ExecuteTemplate(w, "login.tmpl", Data{Title: "Flashlight Login", Error: nil})
	}
}

//Logout logs out a User.
func Logout(w http.ResponseWriter, r *http.Request) {
	mw.SetupSession(w, r)
	defer mw.SaveSession(w, r)
	mw.EndSession()

}

//Register registrates a User if the username isn't taken and the password isn't too short and it matches the password confirmation.
func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("password_confirm")

		var errors []string

		if password != passwordConfirm {
			errors = append(errors, "Passwörter stimmen nicht überein")
		}

		if len(password) < 8 {
			errors = append(errors, "Die Länge des Passworts muss mindestens 8 Zeichen sein")
		}

		//Potentially dangerous due to SQL injection, gotta check!
		if _, err := db.GetUser(username); err == nil {
			errors = append(errors, "Benutzername existiert bereits")
		}

		if len(errors) > 0 {
			mw.Templ.ExecuteTemplate(w, "register.tmpl", Data{Title: "Flashlight Registrieren", Error: errors})
			return
		}

		mw.SetupSession(w, r)
		defer mw.SaveSession(w, r)

		user := db.User{
			Name:     username,
			Type:     "user",
			Password: password,
		}

		db.AddUser(user)

		/*
			body, err := r.GetBody()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			req, err := http.NewRequest("GET", "/", body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			http.Redirect(w, req, "/", http.StatusOK)
		*/
		http.Redirect(w, r, "/", http.StatusOK)
	} else if r.Method == "GET" {
		mw.Templ.ExecuteTemplate(w, "register.tmpl", Data{Title: "Flashlight Registrierung", Error: nil})
	}

}
