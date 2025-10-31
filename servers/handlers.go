package main

import (
	"encoding/json"
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

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type input struct {
		Body string `json:"body"`
	}
	type error struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(req.Body)
	params := input{}
	err := decoder.Decode(&params)
	if err != nil {
		response := error{
			Error: "Something went wrong",
		}
		data, error := json.Marshal(response)
		if error != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(500)
		w.Write(data)
		return
	}

	if len(params.Body) > 140 {
		tooLong := error{
			Error: "Chirp is too long",
		}
		data, error := json.Marshal(tooLong)
		if error != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	type valid struct {
		Valid bool `json:"valid"`
	}
	success := valid{
		Valid: true,
	}
	data, err := json.Marshal(success)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}
