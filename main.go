package main

import (
	"net/http"
)



func main() {

	apiCfg := apiConfig {
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./"))
	handlerFs := http.StripPrefix("/app/",  fs)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handlerFs))

	mux.HandleFunc("GET /api/metrics", apiCfg.printHits)
	mux.HandleFunc("POST /api/reset", apiCfg.resetHits)
	
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


