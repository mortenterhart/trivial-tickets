package server

import (
	"fmt"
	"go-tickets/structs"
	"html/template"
	"log"
	"net/http"
)

// Holds the parsed templates
var tmpl *template.Template

// StartServer gets the parameters for the server and starts it
func StartServer(config *structs.Config) error {

	tmpl = GetTemplates()

	startHandlers()

	return http.ListenAndServe(fmt.Sprintf("%s%d", ":", config.Port), nil)
}

// GetTemplates crawls through the templates folder and reads in all
// present templates.
func GetTemplates() *template.Template {

	// Crawl via relative path, since our current work dir is in cmd/ticketsystem
	t, errTemplates := template.ParseGlob("../../www/templates/*.html")

	if errTemplates != nil {
		log.Fatal("Unable to load the templates: ", errTemplates)
	}

	return t
}

// startHandlers maps all the various handles to the url patterns.
func startHandlers() {
	http.HandleFunc("/", handleIndex)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../www/static"))))
}

// handleIndex handles the traffic for the index.html
func handleIndex(w http.ResponseWriter, r *http.Request) {

	// Render index.html to the browser
	tmpl.Lookup("index.html").ExecuteTemplate(w, "index", nil)
}
