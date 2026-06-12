package main

import (
	"Chirpy/internal/database"
	"net/http"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"log"
)


func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Couldn't open database")
	}

	dbQueries := database.New(db)

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	apiCfg := apiConfig {
		dbQueries: dbQueries,
		platform: platform,
	}
	

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./"))
	handlerFs := http.StripPrefix("/app/",  fs)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handlerFs))

	mux.HandleFunc("GET /admin/metrics", apiCfg.printHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetServer)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChirp)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("GET /api/chirps", apiCfg.returnChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.returnOneChirp)
	mux.HandleFunc("POST /api/login", apiCfg.loginUser)
	
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


