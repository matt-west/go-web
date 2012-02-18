package main

import "fmt"
import "http"

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Awesome! GO Server is UP!")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9980", nil)
}