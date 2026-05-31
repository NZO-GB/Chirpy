package main

import (
	"Chirpy/internal/database"
	"net/http"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)



func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig {
		dbQueries: dbQueries,
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./"))
	handlerFs := http.StripPrefix("/app/",  fs)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handlerFs))

	mux.HandleFunc("GET /admin/metrics", apiCfg.printHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHits)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirp)
	
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	})

	s := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	s.ListenAndServe()
}


