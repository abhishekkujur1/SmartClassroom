package server

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer() {
	port := ":8080"
	fmt.Printf("Starting server on port %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
