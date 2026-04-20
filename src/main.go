package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (not required, only for local development)
	_ = godotenv.Load()

	// Load environment variables
	pushoverAPIToken := os.Getenv("PUSHOVER_API_TOKEN")
	pushoverUserKey := os.Getenv("PUSHOVER_USER_KEY")

	// Set up paths relative to binary location
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	baseDir := filepath.Dir(exePath)

	// Create router
	router := http.NewServeMux()

	// Register handlers
	router.HandleFunc("GET /", handleIndex(baseDir))
	router.HandleFunc("POST /", handleSend(baseDir, pushoverAPIToken, pushoverUserKey))

	// Serve static files
	staticDir := filepath.Join(baseDir, "static")
	fs := http.FileServer(http.Dir(staticDir))
	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
