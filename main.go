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
	r.HandleFunc("/hello", handler).Methods("GET")

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
