package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting server")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", mux)
}
