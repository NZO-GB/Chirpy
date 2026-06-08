package main

import(
	"net/http"
	"encoding/json"
	"fmt"
	"strings"
)

const censor string = "****"

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

func decodeResponse[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var v T

	decoder := json.NewDecoder(r.Body)
	
	err := decoder.Decode(&v)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, err)
		return v, err
	}

	return v, nil
}
	