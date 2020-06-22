package middleware

import (
	"encoding/base64"
	"net/http"
	"os"

	"github.com/gorilla/sessions"

	"github.com/gorilla/securecookie"

	"golang.org/x/crypto/bcrypt"
)

// EnvSessionKey is the String under which we find the Session Key in the Environment Variables
const EnvSessionKey string = "SESSION_KEY"
const SessionSTR string = "session"
const AuthenticatedSTR string = "authenticated"
const UsernameSTR string = "username"

func init() {
	sessionKey := os.Getenv(EnvSessionKey)
	// Create a new Session Key if there isn't one saved yet
	if sessionKey == "" {
		key := string(securecookie.GenerateRandomKey(32))
		os.Setenv(EnvSessionKey, key)
	}
}

// We need a Cookie Store with a private Key. This key should be generated once.
var store = sessions.NewCookieStore([]byte(os.Getenv(EnvSessionKey)))

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := database.GetUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if passwordCorrect(user.Password, password) {
		session, _ := store.Get(r, SessionSTR)

		session.Values[AuthenticatedSTR] = true
		session.Values[UsernameSTR] = username

		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		tmpl.ExecuteTemplate(w, "login.tmpl", nil)
	}
}

func passwordCorrect(passwordHashed string, passwordPlain string) bool {
	passwordDB, _ := base64.StdEncoding.DecodeString(passwordHashed)
	err := bcrypt.CompareHashAndPassword(passwordDB, []byte(passwordPlain))
	if err == nil {
		return true
	}
	return false
}
