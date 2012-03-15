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

func pageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseSetFiles("templates.html")
	// Check that the template file parsed correctly
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}

	slug := r.URL.Path[pagePath:]
	if slug == "" {
		slug = "index"
	}

	c, err := template.ParseFile("pages/" + slug + ".html")
	// Check that the page exists
	if err != nil {
		e, _ := template.ParseFile("errors/404.html")

		w.WriteHeader(http.StatusNotFound)
		e.Execute(w, nil)
		return
	}

	p := findPage(slug)

	t.Execute(w, "Header", p)
	c.Execute(w, nil)
	t.Execute(w, "Footer", nil)
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	assetFile := r.URL.Path[pagePath:]
	http.ServeFile(w, r, assetFile)
}

func findPage(slug string)(page Page) {
	pagesRaw, _ := ioutil.ReadFile("pages/pages.json")

	var pages []Page
	err := json.Unmarshal(pagesRaw, &pages)

	if err != nil {
		// Do Something
	}

	//log.Print(pages)
	for i:=0;i<len(pages);i++ {
		if pages[i].Slug == slug {
			page.Slug = pages[i].Slug
			page.Title = pages[i].Title
			page.Keywords = pages[i].Keywords
			page.Description = pages[i].Description
		}
	}

	return page
}

func main() {
	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/assets/", assetHandler)
	http.ListenAndServe(":9980", nil)
}