package controller

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFiles(
		"templates/footer.tmpl",
		"templates/header.tmpl",
		"templates/footer.tmpl",
		"templates/index_top_logged_in.tmpl",
		"templates/index_top_logged_out.tmpl",
		"templates/nav_index.tmpl",
		"templates/nav_elsewhere.tmpl",
		"templates/index.tmpl",
		"templates/register.tmpl",
		"templates/login.tmpl"))
}

func Preview(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.tmpl", nil)
}
