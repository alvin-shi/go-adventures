package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiState struct {
	fileserverHits atomic.Int32
}

func main() {
	fmt.Println("Starting server")
	state := &apiState{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", state.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("POST /api/reset", state.handlerReset)
	mux.HandleFunc("GET /api/metrics", state.handlerMetrics)
	mux.HandleFunc("GET /api/healthz", handlerHealthCheck)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
