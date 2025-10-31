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
	mux.HandleFunc("GET /api/metrics", state.handlerMetrics)
	mux.HandleFunc("GET /api/healthz", handlerHealthCheck)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("GET /admin/metrics", state.handlerHtmlMetrics)
	mux.HandleFunc("POST /admin/reset", state.handlerReset)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
