package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting server")
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.Handle("/", http.FileServer(http.Dir(".")))
	server.ListenAndServe()
}
