package main

import (
	auth "Chirpy/internal/auth"
	db "Chirpy/internal/database"
	"context"
	"fmt"
	"net/http"
	"time"
	"github.com/alexedwards/argon2id"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	response, err := decodeResponse[UserRequest](w, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	hashedPassword, err := auth.HashPassword(response.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
	}

	userParams := db.CreateUserParams{
		Email: response.Email,
		HashedPassword: hashedPassword,
	}


	user, err := cfg.dbQueries.CreateUser(context.Background(), userParams)
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

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {

	response, err := decodeResponse[UserRequest](w, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	user, err := cfg.dbQueries.ReturnUserByEmail(context.Background(), response.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(response.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	if !match {
		err = fmt.Errorf("Incorrect email or password")
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}

	expiresIn := time.Hour

	token, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return 
	}

	refreshToken := auth.MakeRefreshToken()
	expiresAt :=time.Now().Add(time.Hour * 24 * 60) 

	createTokenArgs := db.CreateTokenParams{
		Token:		token,
		UserID:		user.ID,
		ExpiresAt:	expiresAt,
	}

	cfg.dbQueries.CreateToken(context.Background(), createTokenArgs)

	respondWithJSON(w, http.StatusOK, UserJSON{
			ID:				user.ID,
			Created_at: 	user.CreatedAt,
			Updated_at: 	user.UpdatedAt,
			Email:			user.Email,
			Token:			token,
			Refresh_Token:	refreshToken,
			})
}