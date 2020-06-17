// This is the name of our package
// Everything with this package name can see everything
// else inside the same package, regardless of the file they are in
package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (
	// "fmt" has methods for formatted I/O operations (like printing to the console)
	"fmt"
	"io/ioutil"
	"os"

	// The "net/http" library has methods to implement HTTP clients and servers
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gorilla/sessions"

	"github.com/gorilla/securecookie"
)

// EnvSessionKey is the String under which we find the Session Key in the Environment Variables
const EnvSessionKey string = "SESSION_KEY"

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

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")
	r.HandleFunc("/upload", uploadHandler).Methods("POST")
	r.HandleFunc("/secret", secret).Methods("GET")
	r.HandleFunc("/login", userLogin).Methods("GET")
	r.HandleFunc("/logout", userLogout).Methods("GET")

	// This is the directory we want to publish, in this case,
	// the project root, which is currently our working directory.

	projectRootDir := http.Dir(".")
	staticFileDir := http.Dir("./assets/")

	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/files/" prefix when looking for files.
	// For example, if we type "/files/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for
	// "./files/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/files/", http.FileServer(projectRootDir))

	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/files/", instead of the absolute route itself
	r.PathPrefix("/files/").Handler(staticFileHandler).Methods("GET")

	// Same as above, just for our main pages to be served on the project root

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(staticFileDir))).Methods("GET")
	return r
}

func main() {

	r := newRouter()
	//http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("."))))
	// After defining our server, we finally "listen and serve" on port 8080
	// The second argument is the handler, which we defined earlier
	// and the handler defined above (in "HandleFunc") is used
	http.ListenAndServe(":8080", r)
}

// "handler" is our handler function. It has to follow the function signature of a ResponseWriter and Request type
// as the arguments.
func handler(w http.ResponseWriter, r *http.Request) {
	// For this case, we will always pipe "Hello World" into the response writer
	fmt.Fprintf(w, "Hello World!")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	_, err := store.Get(r, "session-name")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse up to 100 Megabytes of filesize
	r.ParseMultipartForm(100000000)

	file, handler, err := r.FormFile("newImage")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded Filename: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile(os.TempDir(), "upload-*.png")
	if err != nil {
		http.Error(w, "Cannot create temporary File! "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	description := r.FormValue("description")

	fmt.Fprintf(w, "The file can be found under %+v\n", tempFile.Name())

	fmt.Fprintf(w, "The Description for this image is:\n%+v\n", description)

	fmt.Fprintf(w, "Upload successful\n")
}

func secret(w http.ResponseWriter, r *http.Request) {

}

func userLogin(w http.ResponseWriter, r *http.Request) {

}

func userLogout(w http.ResponseWriter, r *http.Request) {

}
