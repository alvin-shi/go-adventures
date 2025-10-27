package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiState struct {
	fileserverHits atomic.Int32
}

func (state *apiState) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (state *apiState) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	hits := fmt.Sprintf("Hits: %v", state.fileserverHits.Load())
	w.Write([]byte(hits))
}
func (state *apiState) handlerReset(w http.ResponseWriter, req *http.Request) {
	state.fileserverHits.Store(0)
}
func handlerHealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte("OK"))
}

func main() {
	fmt.Println("Starting server")
	state := &apiState{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", state.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/reset", state.handlerReset)
	mux.HandleFunc("/metrics", state.handlerMetrics)
	mux.HandleFunc("/healthz", handlerHealthCheck)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
