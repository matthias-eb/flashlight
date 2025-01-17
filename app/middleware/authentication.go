package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"

	"github.com/gorilla/securecookie"

	"golang.org/x/crypto/bcrypt"
)

const envSessionKey string = "SESSION_KEY"
const sessionSTR string = "session"
const authenticatedSTR string = "authenticated"
const usernameSTR string = "username"

func init() {
	sessionKey := os.Getenv(envSessionKey)
	// Create a new Session Key if there isn't one saved yet
	if sessionKey == "" {
		key := string(securecookie.GenerateRandomKey(32))
		os.Setenv(envSessionKey, key)
	}
	cookieStore = sessions.NewCookieStore([]byte(os.Getenv(envSessionKey)))

}

//cookieStore : We need a cookieStore with a private Key. This key should be generated once.
var cookieStore *sessions.CookieStore

var session *sessions.Session

//SetupSession saves the session after everything else has run.
//This is a handler Method. It gets invoked through alice in main.go
func SetupSession(w http.ResponseWriter, r *http.Request) {
	var err error
	session, err = cookieStore.Get(r, sessionSTR)
	if err != nil {
		fmt.Printf("Error getting Session: %+v\n", err.Error())
	}
	if session.IsNew {
		session.Values[authenticatedSTR] = false
		session.Values[usernameSTR] = ""
		SaveSession(w, r)
	}
}

//SaveSession saves the current session.
func SaveSession(w http.ResponseWriter, r *http.Request) {
	session.Save(r, w)
	fmt.Printf("Session saved\n")
}

//AuthenticateUser authenticates a User and saves the cookies, if the password is correct.
func AuthenticateUser(username string, password string, hashedPassword string) (err error) {

	if passwordCorrect(hashedPassword, password) {
		session.Values[authenticatedSTR] = true
		session.Values[usernameSTR] = username
	} else {
		err = errors.New("Passwords didn't Match or User not correct")
		return
	}
	return nil
}

//EndSession unauthenticates the User and removes the Username from the Values.
func EndSession(w http.ResponseWriter, r *http.Request) {
	username := session.Values[usernameSTR]
	fmt.Printf("Logging out user %s\n", username)
	session.Values[usernameSTR] = ""
	session.Values[authenticatedSTR] = false
	SaveSession(w, r)
}

//CheckAuthentication checks if the current User has a valid session and if the session is authenticated. If it is not, then an Error message should be returned and the Starting Page should be redircted to.
//Returns the username and nil if no error occured and an empty username and the error otherwise.
func CheckAuthentication(w http.ResponseWriter, r *http.Request) (string, error) {
	username := session.Values[usernameSTR].(string)
	if !session.Values[authenticatedSTR].(bool) {
		return username, errors.New("User not authenticated")
	}
	fmt.Printf("User %+v is logged in.\n", username)
	return username, nil
}

// HashPassword hashes the given password to save into the database.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(bytes), err
}

func passwordCorrect(passwordHashed string, passwordPlain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(passwordPlain))
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}
