package main

import (
	"sync/atomic"
	"net/http"
	"fmt"
	"encoding/json"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	}

func respondWithError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	msg := fmt.Sprintf(`<html>
  <body>
    <h1>An error has occured:</h1>
    <p>%s</p>
  </body>
</html>`, err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
}

func respondWithJSON(w http.ResponseWriter, json string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) printHits(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	hits := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHits.Store(0)
	w.Write([]byte(`Server hits reset`))
}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {

	var body []byte
	r.Body.Read(body)

	var text string
	err := json.Unmarshal([]byte(body), &text)
	if err != nil {
		respondWithError(w, err)
		return
	}

	
}