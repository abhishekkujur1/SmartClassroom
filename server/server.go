package server

import (
	"fmt"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Smart Attendence Server is running!\n")
}

func StartServer() {
	http.HandleFunc("/health", healthHandler)

	port := ":8080"
	fmt.Printf("Starting server on port %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

