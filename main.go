package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (
	// "fmt" has methods for formatted I/O operations (like printing to the console)

	//Chaining Methods with this Handler

	// Import own controller functions
	"fmt"

	ctr "github.com/matthias-eb/flashlight/app/controller"
	mw "github.com/matthias-eb/flashlight/app/middleware"

	// The "net/http" library has methods to implement HTTP clients and servers
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", mw.SetupMiddleware(ctr.Preview)).Methods("GET")
	r.HandleFunc("/login", ctr.Login).Methods("GET", "POST")
	r.HandleFunc("/register", ctr.Register).Methods("GET", "POST")
	r.HandleFunc("/logout", ctr.Logout).Methods("POST")
	r.HandleFunc("/upload", ctr.UploadImage).Methods("GET", "POST")
	r.HandleFunc("/my-images", ctr.GetImages).Methods("GET")
	r.HandleFunc("/comment", ctr.AddComment).Methods("POST")
	r.HandleFunc("/like", ctr.LikeImage).Methods("POST")
	r.HandleFunc("/deleteImage", ctr.DeleteImage).Methods("POST")

	// This is the directory we want to publish, in this case,
	// the project root, which is currently our working directory.

	projectRootDir := http.Dir(".")
	staticFileDir := http.Dir("./assets/css")
	staticArtefactsDir := http.Dir("./assets/artefacts")

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

	r.PathPrefix("/artefacts/").Handler(http.StripPrefix("/artefacts/", http.FileServer(staticArtefactsDir))).Methods("GET")
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(staticFileDir))).Methods("GET")
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("images/")))).Methods("GET")
	return r
}

func main() {

	r := newRouter()
	//http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("."))))
	// After defining our server, we finally "listen and serve" on port 8080
	// The second argument is the handler, which we defined earlier
	// and the handler defined above (in "HandleFunc") is used
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
