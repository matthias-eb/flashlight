package controller

import (
	"net/http"

	mw "github.com/matthias-eb/flashlight/app/middleware"
)

//Preview responds to a Get Request to Root. It will then show the index Page with the newest Posts of all Users
func Preview(w http.ResponseWriter, r *http.Request) {
	mw.Templ.ExecuteTemplate(w, "index.tmpl", nil)
}
