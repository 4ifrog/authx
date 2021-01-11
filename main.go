package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	port = "8080"
)

func loginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})
}

// Configured as an API
func main() {
	http.Handle("/", loginHandler())

	log.Println("web server running at port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
