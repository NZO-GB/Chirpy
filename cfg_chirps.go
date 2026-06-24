package main

import (
	auth "Chirpy/internal/auth"
	db "Chirpy/internal/database"
	"context"
	"fmt"
	"net/http"
	uuid "github.com/google/uuid"
)

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {

    tokenString, err := auth.GetBearerToken(r.Header)
    if err != nil {
		errReply := fmt.Errorf("Error getting bearing token: %v", err)
        respondWithError(w, http.StatusUnauthorized, errReply)
        return
    }

    userID, err := auth.ValidateJWT(tokenString, cfg.secret)
    if err != nil {
		errReply := fmt.Errorf("Error validating jwt: %v. The jwt being validated was: %v", err, tokenString)
        respondWithError(w, http.StatusUnauthorized, errReply)
        return
    }

	type chirpRequest struct {
        Body string `json:"body"`
    }

    response, err := decodeResponse[chirpRequest](w, r)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err)
        return
    }

    chirpyText, err := validateChirp(response.Body)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err)
        return
    }

    chirpParams := db.CreateChirpParams{
        Body:   chirpyText,
        UserID: userID, 
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
		User_id: 	chirp.UserID,
		}

		cleanChirps[i] = cleanChirp
	}
	respondWithJSON(w, http.StatusOK, cleanChirps)
}

func (cfg *apiConfig) returnOneChirp(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("chirpID")

	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Errorf("Error parsing the uuid"))
		return
	}

	chirp, err := cfg.dbQueries.GetOneChirp(context.Background(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Errorf("Error, chirp not in database"))
		return
	}

	respondWithJSON(w, http.StatusOK, ChirpJSON{
		ID:			chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body: 		chirp.Body,
		User_id: chirp.UserID,
	})
}