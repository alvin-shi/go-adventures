package main

import (
	"fmt"
	"net/http"
)

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

func (state *apiState) handlerHtmlMetrics(w http.ResponseWriter, req *http.Request) {
	hits := fmt.Sprintf("Chirpy has been visited %d times!", state.fileserverHits.Load())
	template := `<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>%s</p>
	</body>
</html>
	`

	html := fmt.Sprintf(template, hits)
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (state *apiState) handlerReset(w http.ResponseWriter, req *http.Request) {
	state.fileserverHits.Store(0)
}
func handlerHealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "plain/text; charset=utf-8")
	w.Write([]byte("OK"))
}
