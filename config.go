package main

import (
	db "Chirpy/internal/database"
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
)



type apiConfig struct {
	fileserverHits	atomic.Int32
	dbQueries		*db.Queries
	platform		string
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

func (cfg *apiConfig) resetServer(w http.ResponseWriter, _ *http.Request) {
	if cfg.platform != "dev" {
		err := fmt.Errorf("You're not an admin")
		respondWithError(w, http.StatusForbidden, err)
	}

	err := cfg.dbQueries.ResetChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
	}

	err = cfg.dbQueries.ResetUsers(context.Background())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
	}

	respondWithJSON(w, http.StatusOK, "Server has been reset")
}

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {

	type chirpRequest struct {
		Body 		string 		`json:"body"`
		User_id		string 		`json:"user_id"`
	}

	response, err := decodeResponse[chirpRequest](w, r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	chirpyText := response.Body

	const maxChirpLength = 140
	if len(chirpyText) > maxChirpLength {
		err := fmt.Errorf("Chirpy is above 140 characters")
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	chirpyText = censorWords(chirpyText, censoringBank)

	chirpParams := db.CreateChirpParams{
		Body: chirpyText,
		UserID: response.User_id,
	}

	chirp, err := cfg.dbQueries.CreateChirp(context.Background(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, ChirpJSON{
		ID:			chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body: 		chirp.Body,
		User_id: chirp.UserID,

	})

}

func (cfg *apiConfig) returnUser(w http.ResponseWriter, r *http.Request) {

	response, err := decodeResponse[UserJSON](w, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(context.Background(), response.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, UserJSON{
    ID:        user.ID,
    Created_at: user.CreatedAt,
    Updated_at: user.UpdatedAt,
    Email:     user.Email,
	})
}

func (cfg *apiConfig) returnChirps(w http.ResponseWriter, _ *http.Request) {


	chirps, err := cfg.dbQueries.GetChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	cleanChirps := make([]ChirpJSON, len(chirps))

	for i, chirp := range chirps {
		cleanChirp := ChirpJSON{
		ID:			chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body: 		chirp.Body,
		User_id: chirp.UserID,
		}
		cleanChirps[i] = cleanChirp
	}

	respondWithJSON(w, http.StatusOK, cleanChirps)
}