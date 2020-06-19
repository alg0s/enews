package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type payload struct {
	text  string
	props string
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	var p payload

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "Payload: ", p)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/handle", handler)
	mux.HandleFunc("/steve", hello)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
