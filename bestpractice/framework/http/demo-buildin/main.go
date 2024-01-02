package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contact", contactHandler)

	port := ":8080"
	fmt.Printf("Start server on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Fprintln(writer, "Welcome to home!")
}

func aboutHandler(writer http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Fprintln(writer, "Welcome to about!")
}

func contactHandler(writer http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Fprintln(writer, "Welcome to contact!")
}
