package main

import "net/http"

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

