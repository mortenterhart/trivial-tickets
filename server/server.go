package server

import (
	"fmt"
	"go-tickets/structs"
	"html/template"
	"log"
	"net/http"
)

// StartServer gets the parameters for the server and starts it
func StartServer(config *structs.Config) error {

	startHandlers()

	return http.ListenAndServe(fmt.Sprintf("%s%d", ":", config.Port), nil)
}

// GetTemplates crawls through the templates folder and reads in all
// present templates.
func GetTemplates() (*template.Template, error) {

	// Crawl via relative path, since our current work dir is in cmd/ticketsystem
	return template.ParseGlob("../../www/templates/*.html")
}

// startHandlers maps all the various handles to the url patterns.
func startHandlers() {
	http.HandleFunc("/", handleIndex)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

// handleIndex handles the traffic for the index.html
func handleIndex(w http.ResponseWriter, r *http.Request) {

	tmpl, errTemplates := GetTemplates()

	if errTemplates != nil {
		log.Fatal("Unable to load the templates: ", errTemplates)
	}

	tmpl.Lookup("index.html").ExecuteTemplate(w, "index", r)
}
