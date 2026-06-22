package main

import (
	auth "Chirpy/internal/auth"
	db "Chirpy/internal/database"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)


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
		return
	}

	err := cfg.dbQueries.ResetChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	err = cfg.dbQueries.ResetUsers(context.Background())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	cfg.fileserverHits.Store(0)

	respondWithJSON(w, http.StatusOK, "Server has been reset")
}

func (cfg *apiConfig) issueRefreshToken(w http.ResponseWriter, r *http.Request) {

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	retrievedToken, err := cfg.dbQueries.RetrieveToken(context.Background(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
	}

	if !retrievedToken.RevokedAt.Valid || retrievedToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, err)
	}

	cfg.dbQueries.CreateToken()

	respondWithJSON(w, http.StatusNoContent, struct{
		token		string
	} {
		token:		auth.MakeRefreshToken(),
	}) 

}

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, _ *http.Request) {

}
