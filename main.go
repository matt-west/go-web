package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type Config struct {
	URL         string
	Title       string
	Description string
	Lang        string
	Webmaster   string
}

type Page struct {
	Slug        string
	Title       string
	Keywords    string
	Description string
}

type Sitemap struct {
	Config *Config
	Pages  []Page
}

const assetPath = len("/")
const pagePath = len("/")

// Config
var config = new(Config)

// Pages
var pages = make(map[string]*Page)
var pagesJSON []Page
var pageTemplates = make(map[string]*template.Template)

// Templates
var layoutTemplates *template.Template
var errorTemplates *template.Template
var sitemapTemplate *template.Template

// Static Assets i.e. Favicons or Humans.txt
var staticAssets = []string{"humans.txt", "favicon.ico"}

// Init Function to Load Template Files and JSON Dict to Cache
func init() {
	log.Println("Loading Config")
	loadConfig()

	log.Println("Loading Templates")
	loadTemplates()

	log.Println("Loading Pages")
	loadPages()
}

// Load the Config File (config/app.json)
func loadConfig() {
	configRaw, _ := ioutil.ReadFile("config/app.json")
	err := json.Unmarshal(configRaw, config)

	if err != nil {
		panic("Could not parse config file!")
	}
}

// Load Pages Dict and Templates
func loadPages() {
	pagesRaw, _ := ioutil.ReadFile("data/pages.json")
	err := json.Unmarshal(pagesRaw, &pagesJSON)
	if err != nil {
		panic("Could not parse Pages JSON!")
	}

	for i := 0; i < len(pagesJSON); i++ {
		pages[pagesJSON[i].Slug] = &pagesJSON[i]
	}

	for _, tmpl := range pages {
		t := template.Must(template.ParseFiles("./pages/" + tmpl.Slug + ".html"))
		pageTemplates[tmpl.Slug] = t
	}
}

// Load Layout and Error Templates
func loadTemplates() {
	layoutTemplates = template.Must(template.ParseFiles("./templates/layouts.html"))
	errorTemplates = template.Must(template.ParseFiles("./templates/errors/404.html", "./templates/errors/505.html"))
	sitemapTemplate = template.Must(template.ParseFiles("./templates/sitemap.xml"))
}

// Page Handler Constructs and Serves Pages
func pageHandler(w http.ResponseWriter, r *http.Request) {
	// Check to see if the request is after a static asset
	for _, asset := range staticAssets {
		if asset == r.URL.Path[1:] {
			http.ServeFile(w, r, asset)
			return
		}
	}

	// Get the page slug, use 'index' if no slug is present
	slug := r.URL.Path[pagePath:]
	if slug == "" {
		slug = "index"
	}

	// Check that the page exists and return 404 if it doesn't
	_, ok := pages[slug]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		errorTemplates.ExecuteTemplate(w, "404", nil)
		return
	}

	// Find the page
	p := pages[slug]

	// Header
	layoutTemplates.ExecuteTemplate(w, "Header", p)

	// Page Template
	err := pageTemplates[slug].Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorTemplates.ExecuteTemplate(w, "505", nil)
		return
	}

	// Footer
	layoutTemplates.ExecuteTemplate(w, "Footer", nil)
}

// Asset Handler Serves CSS, JS and Images
func assetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[assetPath:])
}

// Sitemaps
func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	sitemap := Sitemap{config, pagesJSON}
	sitemapTemplate.Execute(w, sitemap)
}

// Starts Server and Routes Requests
func main() {
	log.Println("Starting: " + config.Title)

	http.HandleFunc("/assets/", assetHandler)
	http.HandleFunc("/sitemap", sitemapHandler)
	http.HandleFunc("/", pageHandler)

	err := http.ListenAndServe(":9981", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
