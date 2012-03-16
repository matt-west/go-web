package main

import (
	"http"
	"template"
	"io/ioutil"
	"json"
)

type Page struct {
	Slug        string
	Title       string
	Keywords    string
	Description string
}

const pagePath = len("/")

var pages = make(map[string]*Page)
var pageTemplates = make(map[string]*template.Template)
var layoutTemplates *template.Set

// Init Function to Load Template Files and JSON Dict to Cache
func init() {
	// Parse Page JSON Dict
	pagesRaw, _ := ioutil.ReadFile("pages/pages.json")
	var pagesJSON []Page
	err := json.Unmarshal(pagesRaw, &pagesJSON)
	if err != nil {
		// Do Something
	}

	// Put Pages into pages map
	for i := 0; i < len(pagesJSON); i++ {
		pages[pagesJSON[i].Slug] = &pagesJSON[i]
	}

	// Parse and Cache Page Templates
	for _, tmpl := range pages {
		t := template.Must(template.ParseFile("./pages/" + tmpl.Slug + ".html"))
		pageTemplates[tmpl.Slug] = t
	}

	// Parse and Cache Layout Templates
	layoutTemplates, err = template.ParseSetFiles("templates.html")
	if err != nil {
		// Do Something
	}
}

// Page Handler Constructs and Serves Pages
func pageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Remove un-neccessary white space from the file

	// Get the page slug, use 'index' if no slug is present
	slug := r.URL.Path[pagePath:]
	if slug == "" {
		slug = "index"
	}

	// Check that the page exists and return 404 if it doesn't
	_, ok := pages[slug]
	if !ok {
		e, _ := template.ParseFile("errors/404.html")

		w.WriteHeader(http.StatusNotFound)
		e.Execute(w, nil)
		return
	}

	// Find the page
	p := pages[slug]

	// Header
	layoutTemplates.Execute(w, "Header", p)
	layoutTemplates.Execute(w, "Nav", nil)

	// Page Template
	err := pageTemplates[slug].Execute(w, nil)
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}

	// Footer
	layoutTemplates.Execute(w, "Footer", nil)
}

// Asset Handler Serves CSS, JS and Images
func assetHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Cache Assets
	assetFile := r.URL.Path[pagePath:]
	http.ServeFile(w, r, assetFile)
}

// Starts Server and Routes Requests
func main() {
	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/assets/", assetHandler)
	http.ListenAndServe(":9980", nil)
}
