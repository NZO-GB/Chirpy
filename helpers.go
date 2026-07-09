package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	db "Chirpy/internal/database"
	"sort"
)

const censor string = "****"
var censoringBank = []string{
"kerfuffle",
"sharbert",
"fornax",
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	msg := fmt.Sprintf(`<html>
  <body>
    <h1>An error has occured:</h1>
    <p>%s</p>
  </body>
</html>`, err)
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	w.Write(data)
}

func censorWords(text string, censoring []string) string {
	censored := make(map[string]struct{})

	for _, word := range censoring {
		censored[strings.ToLower(word)] = struct{}{}
	}

	words := strings.Fields(text)

	for i, w := range words {
		lower := strings.ToLower(w)

		if _, ok := censored[lower]; ok {
			words[i] = censor
		}
	}

	return strings.Join(words, " ")
}

func validateChirp(chirpyText string) (string, error) {
	const maxChirpLength = 140
	if len(chirpyText) > maxChirpLength {
		err := fmt.Errorf("Chirpy is above 140 characters")
		return "", err
	}

	chirpyText = censorWords(chirpyText, censoringBank)

	return chirpyText, nil
}


func decodeResponse[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var v T

	decoder := json.NewDecoder(r.Body)
	
	err := decoder.Decode(&v)
	if err != nil{
		errReply := fmt.Errorf("Error decoding response: %v", err,)
		return v, errReply
	}

	return v, nil
}

func getExpirationTime() time.Time {
	return time.Now().Add(time.Hour * 24 * 60) 
}

func sortChirps(chirps []db.Chirp, sorting string) []db.Chirp {
	if sorting == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	} else {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	}
	return chirps
}