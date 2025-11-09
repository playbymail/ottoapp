// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package rentals implements a server for EmberJS Super Rentals tutorial.
package rentals

import (
	"log"
	"net/http"
	"time"

	"github.com/playbymail/ottoapp/backend/rentals"
)

func runServer() {
	rsv, _ := rentals.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","msg":"pong"}`))
	})
	mux.HandleFunc("/api/rentals.json", rsv.IndexHandler)
	mux.HandleFunc("/api/rentals/{id}", rsv.IdHandler)

	srv := &http.Server{
		Addr:         ":8181",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Println("Rentals API on :8181")
	log.Fatal(srv.ListenAndServe())
}
