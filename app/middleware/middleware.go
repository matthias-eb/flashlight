package middleware

import "net/http"

func callAfter(w http.ResponseWriter, r *http.Request) {

}

func callBefore(w http.ResponseWriter, r *http.Request) {

}

var handlerFunc func(w http.ResponseWriter, r *http.Request)

//SetupMiddleware sets up a Method that calls the handler in between other Methods that are defined in callAfter and callBefore.
func SetupMiddleware(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	handlerFunc = handler
	return Chain
}

//Chain calls all the Methods before the function that are to be called, and the Method, and afterwards all other stuff.
func Chain(w http.ResponseWriter, r *http.Request) {
	callBefore(w, r)
	handlerFunc(w, r)
	callAfter(w, r)
}
