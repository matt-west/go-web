package main

import "fmt"
import "http"

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Awesome! GO Server is UP!")
	matt := Person{"Matt West", 19}
	matt.Hello(w)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9980", nil)
}