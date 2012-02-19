package main

import "http"
import "fmt"

type Person struct {
	Name	string
	Age		int
}

func (p Person) Hello(w http.ResponseWriter) {
	if p.Name != "" {
		fmt.Fprintln(w, "Hello ", p.Name)
	}
}