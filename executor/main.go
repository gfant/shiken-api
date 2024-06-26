package main

import (
	"executor/handler"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/run", handler.RunHandler)
	mux.HandleFunc("/ping", handler.PingHandler)
	corsHandler := setupCORS(mux)

	if err := http.ListenAndServe(":80", corsHandler); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func setupCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass through to the next handler
		next.ServeHTTP(w, r)
	})
}
