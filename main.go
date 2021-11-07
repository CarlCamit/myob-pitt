package main

import (
	"log"
	"net/http"

	"github.com/carlcamit/myob-pitt/github"
	"github.com/carlcamit/myob-pitt/handler"
)

var baseURL = "https://api.github.com"

func main() {
	client := github.NewClient(baseURL)

	router := http.NewServeMux()
	router.Handle("/health", http.HandlerFunc(handler.Health(client)))
	router.Handle("/metadata", http.HandlerFunc(handler.Metadata(client)))
	router.Handle("/", http.HandlerFunc(handler.Root()))

	port := ":8080"
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	log.Printf("starting server on port %s", port)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
