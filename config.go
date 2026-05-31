package main

import (
	"sync/atomic"
	"net/http"
	"fmt"
	"encoding/json"
	"Chirpy/internal/database"
)



type apiConfig struct {
	fileserverHits	atomic.Int32
	dbQueries		*database.Queries
	}

var censoringBank = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
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

	type text struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	txt := text{}
	err := decoder.Decode(&txt)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	chirpyText := txt.Body
	
	if len(txt.Body) > 140 {
		err := fmt.Errorf("Chirpy is above 140 characters")
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	chirpyText = censorWords(chirpyText, censoringBank)

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	payload := returnVals{
		CleanedBody: chirpyText,
	}

	respondWithJSON(w, 200, payload)


}