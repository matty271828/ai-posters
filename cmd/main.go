package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/matty271828/ai-posters/internal/api"
	"github.com/matty271828/ai-posters/internal/server"
)

func main() {
	basepath, err := setupEnv()
	if err != nil {
		log.Fatalf("Failed to set up environment: %v", err)
	}

	myServer, err := server.NewServer(basepath, "8080")
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	// Setup routes
	myServer.AddRoute("/api/generate-image", api.NewAPI().GenerateImageHandler)

	// Start the server
	if err := myServer.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func setupEnv() (string, error) {
	// Attempt to load from .env file, if it exists
	_ = godotenv.Load("local.env")

	// Get the path of the currently running executable
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error determining executable path: %s\n", err)
		return "", fmt.Errorf("Error determining executable path: %w", err)
	}
	return filepath.Dir(execPath), nil
}
