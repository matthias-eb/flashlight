package middleware

import (
	"html/template"
)

//Templ is the Variable under which all templates are saved.
var Templ *template.Template

func init() {
	Templ = template.Must(template.ParseFiles(
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
