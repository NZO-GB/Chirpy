package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./"))
	mux.Handle("/", fs)
	
	s := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	s.ListenAndServe()
}


