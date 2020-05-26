// This is the name of our package
// Everything with this package name can see everything
// else inside the same package, regardless of the file they are in
package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (
	// "fmt" has methods for formatted I/O operations (like printing to the console)
	"fmt"
	// The "net/http" library has methods to implement HTTP clients and servers
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	
	// The "HandleFunc" method accepts a path and a function as arguments
	// (Yes, we can pass functions as arguments, and even trat them like variables in Go)
	// However, the handler function has to have the appropriate signature (as described by the "handler" function below)
	// This is there the router is useful, it allows us to declae methods that
	// this path will be valid for
	r.HandleFunc("/hello", handler).Methods("GET")

	// Declare the static file directory and point it to the
	// directory we just made
	staticFileDirectory := http.Dir("./assets/")

	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for
	// "./assets/assets/index.html", and yield an error
	//staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	staticFileHandler := http.FileServer(staticFileDirectory)

	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	return r
}

func main() {

	r := newRouter()

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