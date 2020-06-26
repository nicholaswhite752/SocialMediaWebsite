package config

import "html/template"

var TPL *template.Template

func init() {
	//Parses all the gohtml files in the templates folder
	//Used to execute HTML code to user from server
	TPL = template.Must(template.ParseGlob("templates/*.gohtml"))
}
