package main

import (
	"http"
	"template"
)

type Page struct {
	Title 		string
}

const pagePath = len("/")

func pageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseSetFiles("templates.html")
	// Check that the template file parsed correctly
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}

	p := &Page{Title: r.URL.Path[pagePath:]}

	if p.Title == "" {
		p.Title = "index"
	}

	c, err := template.ParseFile("pages/" + p.Title + ".html")
	// Check that the page exists
	if err != nil {
		e, _ := template.ParseFile("errors/404.html")

		w.WriteHeader(http.StatusNotFound)
		e.Execute(w, nil)
		return
	}

	t.Execute(w, "Header", p)
	c.Execute(w, nil)
	t.Execute(w, "Footer", nil)
}

func main() {
	http.HandleFunc("/", pageHandler)
	http.ListenAndServe(":9980", nil)
}