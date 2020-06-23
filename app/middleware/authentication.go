package middleware

import (
	"encoding/base64"
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
}

//CookieStore : We need a Cookie CookieStore with a private Key. This key should be generated once.
var CookieStore = sessions.NewCookieStore([]byte(os.Getenv(envSessionKey)))

var session *sessions.Session

//SetupSession saves the session after everything else has run.
//This is a handler Method. It gets invoked through alice in main.go
func SetupSession(w http.ResponseWriter, r *http.Request) {
	var err error
	session, err = CookieStore.Get(r, sessionSTR)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

//SaveSession saves the current session.
func SaveSession(w http.ResponseWriter, r *http.Request) {
	session.Save(r, w)
}

//AuthenticateUser authenticates a User and saves the cookies, if the password is correct.
func AuthenticateUser(username string, password string, hashedPassword string) (err error) {

	if passwordCorrect(hashedPassword, password) {
		session.Values[authenticatedSTR] = true
		session.Values[usernameSTR] = username
	} else {
		err = errors.New("Passwords didn't Match")
		return
	}
	return nil
}

//EndSession unauthenticates the User and removes the Username from the Values.
func EndSession() {
	fmt.Printf("Logging out user %s", session.Values[usernameSTR])
	session.Values[usernameSTR] = ""
	session.Values[authenticatedSTR] = false
}

//CheckAuthentication checks if the current User has a valid session and if the session is authenticated. If it is not, then an Error message should be returned and the Starting Page is opened
func CheckAuthentication(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		if !session.Values[authenticatedSTR].(bool) {
			Templ.ExecuteTemplate(w, "index.tmpl", nil)
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func passwordCorrect(passwordHashed string, passwordPlain string) bool {
	passwordDB, err := base64.StdEncoding.DecodeString(passwordHashed)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(passwordDB, []byte(passwordPlain))
	if err != nil {
		return false
	}
	return true
}
