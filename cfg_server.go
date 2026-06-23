package main

import (
	auth "Chirpy/internal/auth"
	db "Chirpy/internal/database"
	"context"
	"fmt"
	"net/http"
	"time"
	"database/sql"
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

	err = cfg.dbQueries.ResetTokens(context.Background())
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
		errReply := fmt.Errorf("Error getting bearing token: %v", err)
		respondWithError(w, http.StatusInternalServerError, errReply)
		return
	}

	retrievedToken, err := cfg.dbQueries.RetrieveToken(context.Background(), bearerToken)
	if err != nil {
		errReply := fmt.Errorf("Error retrieving token: %v", err)
		respondWithError(w, http.StatusUnauthorized, errReply)
		return
	}

	if retrievedToken.RevokedAt.Valid || retrievedToken.ExpiresAt.Before(time.Now()) {
		errReply := fmt.Errorf("Token revoked or expired: %v", err)
		respondWithError(w, http.StatusUnauthorized, errReply)
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(context.Background(), retrievedToken.Token)
	if err != nil {
		errReply := fmt.Errorf("Error retrieving user from refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, errReply)
		return
	}

	finalToken, err := auth.MakeJWT(user, cfg.secret)
	if err != nil {
		errReply := fmt.Errorf("Error making JWT: %v", err)
		respondWithError(w, http.StatusInternalServerError, errReply)
		return
	}

	respondWithJSON(w, http.StatusOK, struct{
		Token		string		`json:"token"`
	} {
		Token:		finalToken,
	}) 
}

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errReply := fmt.Errorf("Error getting bearing token: %v", err)
		respondWithError(w, http.StatusInternalServerError, errReply)
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(context.Background(), bearerToken)
	if err != nil {
		errReply := fmt.Errorf("Error retrieving user from refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, errReply)
		return
	}

	nullTime := sql.NullTime{Time: time.Now(), Valid: true}

	params := db.RevokeTokenParams{
		UserID: 	user,
		RevokedAt: 	nullTime,
	}

	err = cfg.dbQueries.RevokeToken(context.Background(), params)
	if err != nil {
		errReply := fmt.Errorf("Error revoking token: %v", err)
		respondWithError(w, http.StatusInternalServerError, errReply)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
