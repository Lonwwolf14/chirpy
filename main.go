package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	const port = ":8080"
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	//For index.html with path "/"
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz/", handlerReadiness)
	log.Fatal(srv.ListenAndServe())

}
