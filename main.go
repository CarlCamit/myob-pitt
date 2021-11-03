package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world\r\n")
	}
	router.Handle("/", http.HandlerFunc(handler))

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
