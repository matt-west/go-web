package main

import (
	"http"
	"template"
	"io/ioutil"
	"json"
)

type Page struct {
	Slug				string
	Title 			string
	Keywords		string
	Description	string
}

const pagePath = len("/")
var pages = make(map[string]*Page)
var pageTemplates = make(map[string]*template.Template)

func init() {
	// Parse Page JSON Dict
	pagesRaw, _ := ioutil.ReadFile("pages/pages.json")
	var pagesJSON []Page
	err := json.Unmarshal(pagesRaw, &pagesJSON)
	if err != nil {
		// Do Something
	}

	// Put Pages into pages map
	for i:=0;i<len(pagesJSON);i++ {
		pages[pagesJSON[i].Slug] = &pagesJSON[i]
	}

	// Parse Page Templates
	for _, tmpl := range []string{"index", "about"} {
		t := template.Must(template.ParseFile("./pages/" + tmpl + ".html"))
		pageTemplates[tmpl] = t
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Remove un-neccessary white space from the file

	t, err := template.ParseSetFiles("templates.html")
	// Check that the template file parsed correctly
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}

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
  p := findPage(slug)

	// Header
	t.Execute(w, "Header", p)

	// Page Template
	err = pageTemplates[slug].Execute(w, nil)
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}

	// Footer
	t.Execute(w, "Footer", nil)
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Cache Assets
	assetFile := r.URL.Path[pagePath:]
	http.ServeFile(w, r, assetFile)
}

func findPage(slug string)(page Page) {
	page.Slug = pages[slug].Slug
	page.Title = pages[slug].Title
	page.Keywords = pages[slug].Keywords
	page.Description = pages[slug].Description

	return page
}

func main() {
	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/assets/", assetHandler)
	http.ListenAndServe(":9980", nil)
}